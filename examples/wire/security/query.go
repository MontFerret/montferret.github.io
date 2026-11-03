package security

import (
	"context"
	"errors"

	"github.com/MontFerret/api"
)

// query borrows remote. Its caller closes the runtime after this function returns.
func query(ctx context.Context, remote api.Runtime) (content *api.Content, err error) {
	plan, err := remote.Compile(ctx, api.NewSource("secure.fql", "RETURN @greeting"))
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, plan.Close()) }()

	session, err := plan.NewSession(ctx,
		api.WithParam("greeting", "Hello, Ada!"),
		api.WithOutputContentType("application/json"),
	)
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, session.Close()) }()

	output, err := session.Run(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, output.Close()) }()

	return output.Collect(ctx)
}
