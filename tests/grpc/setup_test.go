package grpc_test

import (
	"context"
	"github.com/bruceneco/go-template/cmd/app"
	"github.com/bruceneco/go-template/config"
	"github.com/bruceneco/go-template/internal/adapters/db/postgres"
	"github.com/bruceneco/go-template/internal/adapters/grpc"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"golang.org/x/sync/semaphore"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	suiteParams struct {
		fx.In
		DB         *postgres.Connection
		GRPCClient *grpc2.ClientConn
	}
)

//nolint:gochecknoglobals // This is a test file
var sem = semaphore.NewWeighted(1)

func setup(t testing.TB, lc any) *suiteParams {
	os.Setenv("GO_ENV", "test")
	os.Setenv("LOG_LEVEL", zerolog.DebugLevel.String())
	ctx := context.Background()
	if err := sem.Acquire(ctx, 1); err != nil {
		t.Fatal(err)
	}
	if lc == nil {
		lc = func() {}
	}
	var params suiteParams
	a := fxtest.New(t,
		app.Inject(),
		fx.Invoke(lc),
		fx.Populate(&params),
		fx.Provide(func(cfg *config.EnvConfig) *grpc2.ClientConn {
			conn, err := grpc2.NewClient(
				"localhost:"+cfg.GRPCPort,
				grpc2.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			return conn
		}),
		fx.Invoke(grpc.Serve),
	)
	a.RequireStart()
	t.Cleanup(func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("LOG_LEVEL")
		a.RequireStop()
		sem.Release(1)
	})
	return &params
}
