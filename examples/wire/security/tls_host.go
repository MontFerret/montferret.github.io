package security

import (
	"crypto/tls"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/server"
	"google.golang.org/grpc/credentials"
)

// The host loads and renews its certificate; Wire borrows the runtime.
func tlsHost(hosted api.Runtime, certificate tls.Certificate) (*server.Server, error) {
	return server.New(hosted, server.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{certificate},
	})))
}
