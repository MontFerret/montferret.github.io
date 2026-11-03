package streaming

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/api"
)

// WriteQuery borrows remote and destination. The invocation context stays live
// through consumption; output cleanup finishes before the function returns.
func WriteQuery(ctx context.Context, remote api.Runtime, source api.Source, destination io.Writer, options ...api.SessionOption) (err error) {
	output, err := remote.Run(ctx, source, options...)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, output.Close()) }()

	return output.Consume(ctx, func(ctx context.Context, chunk []byte) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		n, err := destination.Write(chunk)
		if err == nil && n != len(chunk) {
			err = io.ErrShortWrite
		}

		return err
	})
}
