package security

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/uapi"
	"github.com/MontFerret/wire/client"
	wirev1 "github.com/MontFerret/wire/gen/ferret/wire/v1"
	"github.com/MontFerret/wire/server"
	gooptions "github.com/ziflex/go-options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type (
	certificates struct {
		roots     *x509.CertPool
		server    tls.Certificate
		client    tls.Certificate
		untrusted tls.Certificate
		otherCA   *x509.CertPool
	}
	countedRuntime struct {
		api.Runtime
		plans, sessions, closedPlans, closedSessions, closes atomic.Int32
		verifiedPeers                                        atomic.Int32
	}
	countedPlan struct {
		api.Plan
		owner *countedRuntime
	}
	countedSession struct {
		api.Session
		owner *countedRuntime
	}
	observedListener struct {
		net.Listener
		connections chan *observedConnection
	}
	observedConnection struct {
		net.Conn
		closed chan struct{}
		once   sync.Once
	}
)

func (r *countedRuntime) Version(ctx context.Context) (api.Version, error) {
	if p, ok := peer.FromContext(ctx); ok {
		if info, ok := p.AuthInfo.(credentials.TLSInfo); ok && len(info.State.VerifiedChains) > 0 {
			r.verifiedPeers.Add(1)
		}
	}
	return r.Runtime.Version(ctx)
}

func (r *countedRuntime) Compile(ctx context.Context, source api.Source, options ...api.PlanOption) (api.Plan, error) {
	plan, err := r.Runtime.Compile(ctx, source, options...)
	if err != nil {
		return nil, err
	}
	r.plans.Add(1)
	return &countedPlan{Plan: plan, owner: r}, nil
}

func (r *countedRuntime) Close() error { r.closes.Add(1); return nil }

func (p *countedPlan) NewSession(ctx context.Context, options ...api.SessionOption) (api.Session, error) {
	session, err := p.Plan.NewSession(ctx, options...)
	if err != nil {
		return nil, err
	}
	p.owner.sessions.Add(1)
	return &countedSession{Session: session, owner: p.owner}, nil
}

func (p *countedPlan) Close() error    { p.owner.closedPlans.Add(1); return p.Plan.Close() }
func (s *countedSession) Close() error { s.owner.closedSessions.Add(1); return s.Session.Close() }

func (l *observedListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	observed := &observedConnection{Conn: conn, closed: make(chan struct{})}
	l.connections <- observed
	return observed, nil
}

func (c *observedConnection) Close() error {
	err := c.Conn.Close()
	c.once.Do(func() { close(c.closed) })
	return err
}

func fixture(t *testing.T, build func(api.Runtime) (*server.Server, error)) (string, *countedRuntime, <-chan *observedConnection) {
	t.Helper()
	native, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}
	runtime := &countedRuntime{Runtime: uapi.Wrap(native, "v2.0.0-alpha.62")}
	srv, err := build(runtime)
	if err != nil {
		_ = native.Close()
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		_ = native.Close()
		t.Fatal(err)
	}
	observed := &observedListener{Listener: listener, connections: make(chan *observedConnection, 32)}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(context.Background(), observed) }()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			t.Error(err)
		}
		select {
		case err := <-served:
			if err != nil {
				t.Error(err)
			}
		case <-ctx.Done():
			t.Error("Serve did not settle")
		}
		if runtime.closes.Load() != 0 {
			t.Error("Wire closed the borrowed runtime")
		}
		if runtime.plans.Load() != runtime.closedPlans.Load() || runtime.sessions.Load() != runtime.closedSessions.Load() {
			t.Errorf("resource cleanup: plans %d/%d, sessions %d/%d", runtime.closedPlans.Load(), runtime.plans.Load(), runtime.closedSessions.Load(), runtime.sessions.Load())
		}
		if err := native.Close(); err != nil {
			t.Error(err)
		}
	})
	return listener.Addr().String(), runtime, observed.connections
}

func newCertificates(t *testing.T) certificates {
	t.Helper()
	newCA := func(serial int64) (*x509.Certificate, ed25519.PrivateKey, *x509.CertPool) {
		public, private, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		ca := &x509.Certificate{SerialNumber: big.NewInt(serial), Subject: pkix.Name{CommonName: "example CA"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign}
		der, err := x509.CreateCertificate(rand.Reader, ca, ca, public, private)
		if err != nil {
			t.Fatal(err)
		}
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			t.Fatal(err)
		}
		pool := x509.NewCertPool()
		pool.AddCert(cert)
		return cert, private, pool
	}
	issue := func(ca *x509.Certificate, key ed25519.PrivateKey, serial int64, usage x509.ExtKeyUsage) tls.Certificate {
		public, private, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		cert := &x509.Certificate{SerialNumber: big.NewInt(serial), Subject: pkix.Name{CommonName: "example identity"}, DNSNames: []string{"wire.test"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{usage}}
		der, err := x509.CreateCertificate(rand.Reader, cert, ca, public, key)
		if err != nil {
			t.Fatal(err)
		}
		return tls.Certificate{Certificate: [][]byte{der, ca.Raw}, PrivateKey: private}
	}
	ca, key, roots := newCA(1)
	other, otherKey, otherCA := newCA(2)
	return certificates{roots: roots, otherCA: otherCA, server: issue(ca, key, 3, x509.ExtKeyUsageServerAuth), client: issue(ca, key, 4, x509.ExtKeyUsageClientAuth), untrusted: issue(other, otherKey, 5, x509.ExtKeyUsageClientAuth)}
}

func checkQuery(t *testing.T, ctx context.Context, remote api.Runtime) {
	t.Helper()
	output, err := query(ctx, remote)
	if err != nil {
		t.Fatal(err)
	}
	if output == nil || output.Metadata.ContentType != "application/json" || string(output.Data) != `"Hello, Ada!"` {
		t.Fatalf("output = %#v", output)
	}
}

func TestTLS(t *testing.T) {
	certs := newCertificates(t)
	address, _, _ := fixture(t, func(r api.Runtime) (*server.Server, error) { return tlsHost(r, certs.server) })
	for _, test := range []struct {
		name     string
		connect  func(context.Context) (api.Runtime, error)
		accepted bool
	}{
		{"trusted identity", func(ctx context.Context) (api.Runtime, error) {
			return tlsClient(ctx, address, certs.roots, "wire.test")
		}, true},
		{"untrusted server", func(ctx context.Context) (api.Runtime, error) {
			return tlsClient(ctx, address, certs.otherCA, "wire.test")
		}, false},
		{"wrong identity", func(ctx context.Context) (api.Runtime, error) {
			return tlsClient(ctx, address, certs.roots, "other.test")
		}, false},
		{"plaintext", func(ctx context.Context) (api.Runtime, error) { return client.New(ctx, address, client.WithInsecure()) }, false},
		{"default system trust", func(ctx context.Context) (api.Runtime, error) { return client.New(ctx, address) }, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			budget := 500 * time.Millisecond
			if test.accepted {
				budget = 3 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.Background(), budget)
			defer cancel()
			remote, err := test.connect(ctx)
			if !test.accepted {
				if err == nil || remote != nil {
					t.Fatalf("rejected startup = %v, %v", remote, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := remote.Close(); err != nil {
					t.Error(err)
				}
			}()
			checkQuery(t, ctx, remote)
		})
	}
}

func TestMTLS(t *testing.T) {
	certs := newCertificates(t)
	address, hosted, _ := fixture(t, func(r api.Runtime) (*server.Server, error) { return mtlsHost(r, certs.server, certs.roots) })
	for _, test := range []struct {
		name        string
		certificate *tls.Certificate
	}{
		{"trusted client", &certs.client}, {"missing client", nil}, {"untrusted client", &certs.untrusted},
	} {
		t.Run(test.name, func(t *testing.T) {
			budget := 500 * time.Millisecond
			if test.name == "trusted client" {
				budget = 3 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.Background(), budget)
			defer cancel()
			var remote api.Runtime
			var err error
			if test.certificate == nil {
				remote, err = tlsClient(ctx, address, certs.roots, "wire.test")
			} else {
				remote, err = mtlsClient(ctx, address, certs.roots, "wire.test", *test.certificate)
			}
			if test.name != "trusted client" {
				if err == nil || remote != nil {
					t.Fatalf("rejected startup = %v, %v", remote, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := remote.Close(); err != nil {
					t.Error(err)
				}
			}()
			checkQuery(t, ctx, remote)
			if hosted.verifiedPeers.Load() == 0 {
				t.Fatal("hosted handshake did not see a verified TLS client")
			}
		})
	}
}

func TestTokens(t *testing.T) {
	certs := newCertificates(t)
	var methodsMu sync.Mutex
	methods := make(map[string]int)
	address, _, _ := fixture(t, func(r api.Runtime) (*server.Server, error) {
		return tokenHost(r, certs.server, func(ctx context.Context, token string) error {
			if token != "accepted" {
				return errors.New("private rejection reason")
			}
			method, _ := grpc.Method(ctx)
			methodsMu.Lock()
			methods[method]++
			methodsMu.Unlock()
			return nil
		})
	})
	for _, token := range []string{"", "invalid", "accepted"} {
		t.Run("startup/"+token, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			remote, err := tokenClient(ctx, address, certs.roots, "wire.test", func(context.Context) (string, error) { return token, nil })
			if token != "accepted" {
				if status.Code(err) != codes.Unauthenticated || remote != nil || strings.Contains(err.Error(), "private rejection") {
					t.Fatalf("startup = %v, %v", remote, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			checkQuery(t, ctx, remote)
			if err := remote.Close(); err != nil {
				t.Error(err)
			}
		})
	}
	t.Run("missing startup token", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		remote, err := tlsClient(ctx, address, certs.roots, "wire.test")
		if remote != nil || status.Code(err) != codes.Unauthenticated {
			t.Fatalf("startup = %v, %v", remote, err)
		}
	})

	// Borrow a channel to vary credentials for subsequent unary and stream RPCs.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12, RootCAs: certs.roots, ServerName: "wire.test"})))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	for _, test := range []struct {
		name   string
		token  string
		attach bool
	}{{"missing", "", false}, {"invalid", "invalid", true}} {
		t.Run("operations/"+test.name, func(t *testing.T) {
			var options []grpc.CallOption
			if test.attach {
				options = append(options, grpc.PerRPCCredentials(tokenSource(func(context.Context) (string, error) { return test.token, nil })))
			}
			_, err := wirev1.NewPlanServiceClient(conn).Compile(ctx, &wirev1.CompileRequest{}, options...)
			if status.Code(err) != codes.Unauthenticated {
				t.Fatalf("unary rejection = %v", err)
			}
			stream, err := wirev1.NewExecutionServiceClient(conn).WatchExecution(ctx, &wirev1.WatchExecutionRequest{}, options...)
			if err == nil {
				_, err = stream.Recv()
			}
			if status.Code(err) != codes.Unauthenticated {
				t.Fatalf("stream rejection = %v", err)
			}
		})
	}

	// A refreshing connection-level source also reaches later operations and cleanup.
	var current atomic.Value
	current.Store("accepted")
	remote, err := tokenClient(ctx, address, certs.roots, "wire.test", func(context.Context) (string, error) { return current.Load().(string), nil })
	if err != nil {
		t.Fatal(err)
	}
	current.Store("invalid")
	if _, err := remote.Compile(ctx, api.NewSource("rejected.fql", "RETURN 1")); status.Code(err) != codes.Unauthenticated {
		t.Fatalf("later credentials = %v", err)
	}
	current.Store("accepted")
	checkQuery(t, ctx, remote)
	if err := remote.Close(); err != nil {
		t.Error(err)
	}
	methodsMu.Lock()
	defer methodsMu.Unlock()
	for _, method := range []string{
		"/ferret.wire.v1.RuntimeService/Connect",
		"/ferret.wire.v1.PlanService/Compile",
		"/ferret.wire.v1.SessionService/CreateSession",
		"/ferret.wire.v1.ExecutionService/RunSession",
		"/ferret.wire.v1.ExecutionService/ReadExecutionOutput",
		"/ferret.wire.v1.ExecutionService/ReleaseExecution",
		"/ferret.wire.v1.SessionService/ReleaseSession",
		"/ferret.wire.v1.PlanService/ReleasePlan",
		"/ferret.wire.v1.RuntimeService/CloseConnection",
	} {
		if methods[method] == 0 {
			t.Errorf("accepted credentials did not reach %s", method)
		}
	}
}

func TestChannelOwnership(t *testing.T) {
	certs := newCertificates(t)
	address, _, accepted := fixture(t, func(r api.Runtime) (*server.Server, error) { return tlsHost(r, certs.server) })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	startup, cancelStartup := context.WithCancel(ctx)
	remote, err := tlsClient(startup, address, certs.roots, "wire.test")
	if err != nil {
		t.Fatal(err)
	}
	cancelStartup()
	var owned *observedConnection
	select {
	case owned = <-accepted:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	plan, err := remote.Compile(ctx, api.NewSource("retained.fql", "RETURN 42"))
	if err != nil {
		t.Fatal(err)
	}
	if err := remote.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-owned.closed:
		t.Fatal("owned channel closed before its plan")
	default:
	}
	session, err := plan.NewSession(ctx, api.WithOutputContentType("application/json"))
	if err != nil {
		t.Fatal(err)
	}
	output, err := session.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	content, collectErr := output.Collect(ctx)
	if err := errors.Join(collectErr, output.Close()); err != nil {
		t.Fatal(err)
	}
	if content == nil || string(content.Data) != "42" {
		t.Fatalf("retained session = %#v", content)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if err := plan.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-owned.closed:
	case <-ctx.Done():
		t.Fatal("owned channel was not released")
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12, RootCAs: certs.roots, ServerName: "wire.test"})))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	borrowed, err := client.From(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	checkQuery(t, ctx, borrowed)
	if err := borrowed.Close(); err != nil {
		t.Fatal(err)
	}
	if conn.GetState() == connectivity.Shutdown {
		t.Fatal("Wire closed the borrowed channel")
	}
	second, err := client.From(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	checkQuery(t, ctx, second)
	if err := second.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestValidationExample(t *testing.T) {
	native, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := native.Close(); err != nil {
			t.Error(err)
		}
	}()
	srv, err := validatedHost(uapi.Wrap(native, "v2.0.0-alpha.62"), server.WithLimits(server.Limits{}), nil, server.WithUnaryInterceptors(nil, nil))
	var invalid gooptions.ValidationError
	if srv != nil || !errors.As(err, &invalid) || invalid.Field != "limits" || !invalid.OmitValue || invalid.Value != "" {
		t.Fatalf("validation = %v, %#v", err, invalid)
	}
}

func TestTLSDoesNotFallBack(t *testing.T) {
	address, _, _ := fixture(t, func(r api.Runtime) (*server.Server, error) { return server.New(r) })
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	remote, err := client.New(ctx, address)
	if err == nil || remote != nil {
		t.Fatalf("default TLS reached plaintext host: %v, %v", remote, err)
	}
	secureCtx, secureCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer secureCancel()
	remote, err = client.New(secureCtx, address, client.WithInsecure(), client.WithPerRPCCredentials(tokenSource(func(context.Context) (string, error) { return "secret", nil })))
	if err == nil || remote != nil {
		t.Fatalf("secure credentials reached plaintext transport: %v, %v", remote, err)
	}
}
