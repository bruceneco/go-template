package grpc

import (
	"context"
	"fmt"
	"github.com/bruceneco/go-template/config"
	"net"

	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	Server interface {
		Register(grpcServer *grpc.Server)
	}
	Params struct {
		fx.In
		LC       fx.Lifecycle
		Cfg      *config.EnvConfig
		Servers  []Server `group:"servers"`
		Validate *validator.Validate
	}
)

func NewGRPCServer(p Params) *grpc.Server {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithLevels(logging.DefaultClientCodeToLevel),
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(), opts...),
			ValidateInterceptor(p.Validate),
		),
	)

	if p.Cfg.GoEnv == config.EnvTypeDevelopment {
		reflection.Register(grpcServer)
	}

	for _, s := range p.Servers {
		s.Register(grpcServer)
	}
	return grpcServer
}

func Serve(lc fx.Lifecycle, grpcServer *grpc.Server, cfg *config.EnvConfig) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			lis, err := net.Listen("tcp", "localhost:"+cfg.GRPCPort)
			if err != nil {
				return err
			}

			log.Info().Msgf("GRPC server is running on %s", lis.Addr().String())
			go func() {
				err = grpcServer.Serve(lis)
				if err != nil {
					log.Fatal().Err(err).Msg("GRPC server failed")
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("gracefully stopping GRPC server")

			grpcServer.GracefulStop()

			return nil
		},
	})
	return nil
}

func InterceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		l := log.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
