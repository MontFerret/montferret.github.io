package security

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/client"
	"google.golang.org/grpc/credentials"
)

// roots trusts the issuing CA; serverName must match a certificate DNS identity.
// The returned runtime owns its channel. Close its descendants, then the runtime.
func tlsClient(ctx context.Context, target string, roots *x509.CertPool, serverName string) (api.Runtime, error) {
	return client.New(ctx, target, client.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    roots,
		ServerName: serverName,
	})))
}
