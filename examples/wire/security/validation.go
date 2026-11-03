package security

import (
	"errors"
	"fmt"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/server"
	gooptions "github.com/ziflex/go-options"
)

func validatedHost(hosted api.Runtime, options ...server.Option) (*server.Server, error) {
	srv, err := server.New(hosted, options...)
	if err != nil {
		var invalid gooptions.ValidationError
		if errors.As(err, &invalid) {
			// As finds one failure. Preserve the full joined error for inspection.
			return nil, fmt.Errorf("invalid Wire configuration: %w", err)
		}
		return nil, err
	}
	return srv, nil
}
