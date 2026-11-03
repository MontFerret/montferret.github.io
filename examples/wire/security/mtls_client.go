package security

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/client"
	"google.golang.org/grpc/credentials"
)

func mtlsClient(ctx context.Context, target string, roots *x509.CertPool, serverName string, certificate tls.Certificate) (api.Runtime, error) {
	return client.New(ctx, target, client.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      roots,
		ServerName:   serverName,
		Certificates: []tls.Certificate{certificate},
	})))
}
