// Todo API provides a RESTful API for managing todo items
//
//	@title			Todo API
//	@version		1.0
//	@description	A REST API for managing todo items with clean architecture
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.email	support@example.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:5949
//	@BasePath	/api/v1
//
//	@schemes	http https
//
//	@produce	json
//	@consumes	json
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	_ "github.com/NoroSaroyan/go-rest-api-example/docs" // Import generated docs
	"github.com/NoroSaroyan/go-rest-api-example/internal/app"
	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
	"github.com/NoroSaroyan/go-rest-api-example/internal/shutdown"
)

const (
	// shutdownTimeout defines how long to wait for graceful shutdown
	shutdownTimeout = 5 * time.Second
)

type options struct {
	Log logger.Options `group:"log" namespace:"log" env-namespace:"LOG" description:"Log level"`
}

func main() {
	opts := options{}
	parser := flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Printf("help called: %v\n", err)
			return
		}
		log.Fatal("couldn't parse flags", err)
	}

	ctx := context.Background()
	ctx = logger.Inject(ctx, logger.NewLogger(opts.Log))

	exitCode := runMain(ctx, opts)
	os.Exit(exitCode)
}

func runMain(ctx context.Context, opts options) int {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := logger.FromContext(ctx)
	g, ctx := errgroup.WithContext(ctx)

	application, err := app.New(ctx)
	if err != nil {
		log.Error("failed to create application", zap.Error(err))
		return 1
	}

	g.Go(func() error { return application.Run(ctx) })

	g.Go(func() error { return shutdown.Receive(ctx) })

	errg := g.Wait()

	if errors.Is(errg, shutdown.ErrGracefullyShutdown) {
		log.Info("Server is gracefully shut down")
		return 0
	}
	if errg != nil {
		log.Error("Unexpected error:", zap.Error(errg))
		return 9
	}

	return 0
}
