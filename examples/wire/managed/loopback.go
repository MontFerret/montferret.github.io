package managed

import (
	"context"
	"errors"
	"time"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/server"
)

// ServeLoopback borrows the host's configured runtime. The host closes it only
// after this function returns, including when the managed wait times out.
func ServeLoopback(ctx context.Context, hosted api.Runtime) (err error) {
	srv, err := server.New(hosted)
	if err != nil {
		return err
	}
	defer func() {
		// A timed-out Run can leave hosted cleanup in progress.
		err = errors.Join(err, srv.Shutdown(context.Background()))
	}()

	// Plaintext loopback is an explicit local-development policy.
	return srv.Run(ctx, "127.0.0.1:50051", server.WithShutdownTimeout(30*time.Second))
}
