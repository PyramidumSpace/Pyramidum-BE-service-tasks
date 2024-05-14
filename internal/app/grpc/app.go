package app

import (
	"context"
	"fmt"
	pgconnection "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/database/connection/postgres"
	pgrepository "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/database/repository/postgres"
	apiserver "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	port       int
	grpcServer *grpc.Server
}

func NewApp(log *slog.Logger, port int, database *pgconnection.Database) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall,
			logging.FinishCall,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		),
	)

	r := pgrepository.NewRepository(database.DB())

	apiserver.RegisterServer(grpcServer, log, r)

	return &App{
		port:       port,
		grpcServer: grpcServer,
	}
}

func InterceptorLogger(log *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) Run() error {
	const op = "app.grpc.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
}
