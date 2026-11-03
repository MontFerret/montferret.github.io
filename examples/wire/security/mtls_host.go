package security

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/server"
	"google.golang.org/grpc/credentials"
)

func mtlsHost(hosted api.Runtime, certificate tls.Certificate, clientCAs *x509.CertPool) (*server.Server, error) {
	return server.New(hosted, server.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	})))
}
