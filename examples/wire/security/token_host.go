package security

import (
	"context"
	"crypto/tls"
	"strings"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// validateToken is the host's policy: verify signature, audience, expiry,
// revocation, and permissions as appropriate. Never log the supplied token.
func tokenHost(hosted api.Runtime, certificate tls.Certificate, validateToken func(context.Context, string) error) (*server.Server, error) {
	authenticate := func(ctx context.Context) error {
		values := metadata.ValueFromIncomingContext(ctx, "authorization")
		if len(values) != 1 {
			return status.Error(codes.Unauthenticated, "authentication required")
		}
		token, ok := strings.CutPrefix(values[0], "Bearer ")
		if !ok || token == "" || validateToken(ctx, token) != nil {
			return status.Error(codes.Unauthenticated, "authentication required")
		}
		return nil
	}
	unary := func(ctx context.Context, request any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if err := authenticate(ctx); err != nil {
			return nil, err
		}
		return next(ctx, request)
	}
	stream := func(host any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, next grpc.StreamHandler) error {
		if err := authenticate(stream.Context()); err != nil {
			return err
		}
		return next(host, stream)
	}
	return server.New(hosted,
		server.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{certificate},
		})),
		server.WithUnaryInterceptors(unary),
		server.WithStreamInterceptors(stream),
	)
}
