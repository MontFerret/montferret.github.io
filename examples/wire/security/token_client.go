package security

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/client"
	"google.golang.org/grpc/credentials"
)

// tokenSource may refresh credentials. It also runs for later resource cleanup.
type tokenSource func(context.Context) (string, error)

func (source tokenSource) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	token, err := source(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]string{"authorization": "Bearer " + token}, nil
}

func (tokenSource) RequireTransportSecurity() bool { return true }

func tokenClient(ctx context.Context, target string, roots *x509.CertPool, serverName string, source tokenSource) (api.Runtime, error) {
	return client.New(ctx, target,
		client.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    roots,
			ServerName: serverName,
		})),
		// Connection-level credentials cover the handshake, operations, and cleanup.
		client.WithPerRPCCredentials(source),
	)
}
