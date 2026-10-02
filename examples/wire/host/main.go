package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/uapi"
	"github.com/MontFerret/wire/server"
)

func main() {
	address := flag.String("listen", "127.0.0.1:50051", "local development endpoint")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, *address); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, address string) (resultErr error) {
	native, err := ferret.New(ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Namespace("app").Function().A1().Add("greet", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			name, err := runtime.CastArg[runtime.String](arg, 0)
			if err != nil {
				return nil, err
			}
			return runtime.NewString("Hello, " + string(name) + "!"), nil
		})
	}))
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, native.Close()) }()

	// Wrap and Wire borrow this engine; the host remains its owner.
	hosted := uapi.Wrap(native, "v2.0.0-alpha.57")
	wireServer, err := server.NewServer(hosted)
	if err != nil {
		return err
	}

	// Plaintext transport is for local development only.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer func() {
		if err := listener.Close(); !errors.Is(err, net.ErrClosed) {
			resultErr = errors.Join(resultErr, err)
		}
	}()

	served := make(chan error, 1)
	go func() { served <- wireServer.Serve(context.Background(), listener) }()
	serveReturned := false
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resultErr = errors.Join(resultErr, wireServer.Shutdown(shutdownCtx))
		if !serveReturned {
			select {
			case err := <-served:
				resultErr = errors.Join(resultErr, err)
			case <-shutdownCtx.Done():
				resultErr = errors.Join(resultErr, shutdownCtx.Err())
			}
		}
	}()

	if _, err := fmt.Printf("Listening on %s\n", listener.Addr()); err != nil {
		return err
	}
	select {
	case resultErr = <-served:
		serveReturned = true
	case <-ctx.Done():
		// Shutdown settles Wire-owned work before the host closes its engine.
	}
	return resultErr
}
