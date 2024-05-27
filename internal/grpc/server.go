package grpc

import (
	"context"
	repository "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/database/repository/postgres"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/create"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/get"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/getbyid"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/update"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc"
	"log/slog"
)

type ServerAPI struct {
	proto.UnimplementedTasksServiceServer
	createHandlerFunc      create.HandlerFunc
	getHandlerFunc         get.HandlerFunc
	getByUserIdHandlerFunc getbyid.HandlerFunc
	updateHandlerFunc      update.HandlerFunc
}

func RegisterServer(gRPC *grpc.Server, log *slog.Logger, r *repository.Repository) {
	proto.RegisterTasksServiceServer(
		gRPC,
		&ServerAPI{
			createHandlerFunc:      create.MakeCreateHandler(log, r),
			getHandlerFunc:         get.MakeGetHandler(log, r),
			getByUserIdHandlerFunc: getbyid.MakeGetByUserIdHandler(log, r),
			updateHandlerFunc:      update.MakeUpdateHandler(log, r),
		},
	)
}

func (s *ServerAPI) Create(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
	return s.createHandlerFunc(ctx, req)
}

func (s *ServerAPI) Update(ctx context.Context, req *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	return s.updateHandlerFunc(ctx, req)
}

func (s *ServerAPI) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	return s.getHandlerFunc(ctx, req)
}

func (s *ServerAPI) GetByUserID(ctx context.Context, req *proto.GetByUserIDRequest) (*proto.GetByUserIDResponse, error) {
	return s.getByUserIdHandlerFunc(ctx, req)
}
