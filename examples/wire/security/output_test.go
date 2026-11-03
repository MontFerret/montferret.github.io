package security

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MontFerret/api"
	"github.com/MontFerret/montferret.github.io/examples/wire/streaming"
	"github.com/MontFerret/wire/pkg/failure"
	"github.com/MontFerret/wire/server"
)

type (
	lateFailureRuntime struct {
		api.Runtime
		closedOutputs atomic.Int32
	}

	lateFailureOutput struct {
		api.Output
		owner *lateFailureRuntime
	}

	payloadWriter struct {
		total int
		calls int
	}

	writeFunc func([]byte) (int, error)
)

func (r *lateFailureRuntime) Run(ctx context.Context, source api.Source, options ...api.SessionOption) (api.Output, error) {
	output, err := r.Runtime.Run(ctx, source, options...)
	if err != nil {
		return nil, err
	}

	return &lateFailureOutput{Output: output, owner: r}, nil
}

func (o *lateFailureOutput) Consume(ctx context.Context, consumer api.Consumer) error {
	// Model an encoder that delivers a prefix before discovering a terminal error.
	return o.Output.Consume(ctx, func(ctx context.Context, chunk []byte) error {
		var err error
		if len(chunk) > 0 {
			err = consumer(ctx, chunk[:1])
		}

		return errors.Join(err, errors.New("private encoder failure"))
	})
}

func (o *lateFailureOutput) Close() error {
	o.owner.closedOutputs.Add(1)

	return o.Output.Close()
}

func (w *payloadWriter) Write(data []byte) (int, error) {
	for i, value := range data {
		position := w.total + i
		want := byte('x')
		if position == 0 || position == (8<<20)+1 {
			want = '"'
		}

		if value != want {
			return 0, errors.New("unexpected encoded byte")
		}
	}

	w.total += len(data)
	w.calls++

	return len(data), nil
}

func (w writeFunc) Write(data []byte) (int, error) {
	return w(data)
}

func TestQueryReturnsDetachedContent(t *testing.T) {
	certs := newCertificates(t)
	address, hosted, _ := fixture(t, func(r api.Runtime) (*server.Server, error) { return tlsHost(r, certs.server) })
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	remote, err := tlsClient(ctx, address, certs.roots, "wire.test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := remote.Close(); err != nil {
			t.Error(err)
		}
	})

	content, err := query(ctx, remote)
	if err != nil {
		t.Fatal(err)
	}
	if hosted.closedPlans.Load() != 1 || hosted.closedSessions.Load() != 1 {
		t.Fatal("query returned before its plan and session were closed")
	}
	if err := remote.Close(); err != nil {
		t.Fatal(err)
	}
	cancel()

	if content == nil || content.Metadata.ContentType != "application/json" || string(content.Data) != `"Hello, Ada!"` {
		t.Fatalf("detached content after cleanup = %#v", content)
	}
}

func TestConsumeMultipleMessages(t *testing.T) {
	certs := newCertificates(t)
	address, _, _ := fixture(t, func(r api.Runtime) (*server.Server, error) { return tlsHost(r, certs.server) })
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	remote, err := tlsClient(ctx, address, certs.roots, "wire.test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := remote.Close(); err != nil {
			t.Error(err)
		}
	}()

	// The small request produces more than two default 4 MiB message budgets.
	var destination payloadWriter
	err = streaming.WriteQuery(ctx, remote,
		api.NewAnonymousSource(`RETURN REPEAT("x", 8388608)`),
		&destination, api.WithOutputContentType("application/json"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if destination.total != (8<<20)+2 || destination.calls < 3 {
		t.Fatalf("streamed bytes=%d callbacks=%d", destination.total, destination.calls)
	}
}

func TestTerminalConsumptionFailure(t *testing.T) {
	for _, mode := range []string{"collect", "consume"} {
		t.Run(mode, func(t *testing.T) {
			certs := newCertificates(t)
			var hosted *lateFailureRuntime
			address, _, _ := fixture(t, func(r api.Runtime) (*server.Server, error) {
				hosted = &lateFailureRuntime{Runtime: r}
				return tlsHost(hosted, certs.server)
			})
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			remote, err := tlsClient(ctx, address, certs.roots, "wire.test")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := remote.Close(); err != nil {
					t.Error(err)
				}
			}()

			var received string
			source := api.NewAnonymousSource("RETURN 42")
			if mode == "collect" {
				output, runErr := remote.Run(ctx, source)
				if runErr != nil {
					t.Fatalf("expected successful preparation: %v", runErr)
				}
				content, collectErr := output.Collect(ctx)
				err = errors.Join(collectErr, output.Close())
				if content != nil {
					received = string(content.Data)
				}
			} else {
				var destination bytes.Buffer
				err = streaming.WriteQuery(ctx, remote, source, &destination)
				received = destination.String()
			}

			var terminal *failure.Failure
			if !errors.As(err, &terminal) || terminal.Category != failure.CategoryExecution {
				t.Fatalf("terminal consumption error = %v", err)
			}
			if strings.Contains(err.Error(), "private encoder failure") {
				t.Fatal("host's private failure crossed the boundary")
			}
			if received != "4" || hosted.closedOutputs.Load() != 1 {
				t.Fatalf("received=%q closed outputs=%d", received, hosted.closedOutputs.Load())
			}
		})
	}
}

func TestConsumeWriterFailures(t *testing.T) {
	rejected := errors.New("destination rejected output")
	for _, test := range []struct {
		name   string
		writer io.Writer
		want   error
	}{
		{"writer error", writeFunc(func([]byte) (int, error) { return 0, rejected }), rejected},
		{"short write", writeFunc(func(data []byte) (int, error) { return len(data) - 1, nil }), io.ErrShortWrite},
	} {
		t.Run(test.name, func(t *testing.T) {
			certs := newCertificates(t)
			address, _, _ := fixture(t, func(r api.Runtime) (*server.Server, error) { return tlsHost(r, certs.server) })
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			remote, err := tlsClient(ctx, address, certs.roots, "wire.test")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := remote.Close(); err != nil {
					t.Error(err)
				}
			}()

			err = streaming.WriteQuery(ctx, remote, api.NewAnonymousSource("RETURN 42"), test.writer)
			if !errors.Is(err, test.want) {
				t.Fatalf("writer failure = %v, want %v", err, test.want)
			}
			// Failed consumption has settled its execution; the runtime is reusable.
			checkQuery(t, ctx, remote)
		})
	}
}
