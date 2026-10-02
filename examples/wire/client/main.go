package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MontFerret/api"
	"github.com/MontFerret/wire/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := flag.String("address", "127.0.0.1:50051", "host endpoint")
	name := flag.String("name", "Ada", "greeting parameter")
	timeout := flag.Duration("timeout", 10*time.Second, "maximum query duration")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()
	if err := run(ctx, *address, *name); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, address, name string) (resultErr error) {
	// Plaintext transport is for local development only.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, conn.Close()) }()

	remote, err := client.New(ctx, conn)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, remote.Close()) }()

	plan, err := remote.Compile(ctx, api.NewSource("greeting.fql", "return app::greet(@name)"))
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, plan.Close()) }()

	session, err := plan.NewSession(ctx,
		api.WithParam("name", name),
		api.WithOutputContentType("application/json"),
	)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, session.Close()) }()

	output, err := session.Run(ctx)
	// Available output can accompany an execution or cleanup error.
	if output != nil {
		_, writeErr := fmt.Printf("Content-Type: %s\n%s\n", output.ContentType, output.Content)
		err = errors.Join(err, writeErr)
	}
	return err
}
