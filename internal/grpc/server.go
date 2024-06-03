package grpc

import (
	"context"
	"log/slog"

	repository "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/database/repository/postgres"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/create"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/patch"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/task"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/tasks"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/grpc/handlers/update"
	proto "github.com/pyramidum-space/protos/gen/go/tasks"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	proto.UnimplementedTasksServiceServer
	createHandlerFunc create.HandlerFunc
	taskHandlerFunc   task.HandlerFunc
	tasksHandlerFunc  tasks.HandlerFunc
	updateHandlerFunc update.HandlerFunc
	patchHandlerFunc  patch.HandlerFunc
}

func RegisterServer(gRPC *grpc.Server, log *slog.Logger, r *repository.Repository) {
	proto.RegisterTasksServiceServer(
		gRPC,
		&ServerAPI{
			createHandlerFunc: create.MakeCreateHandler(log, r),
			taskHandlerFunc:   task.MakeTaskHandler(log, r),
			tasksHandlerFunc:  tasks.MakeTasksHandler(log, r),
			updateHandlerFunc: update.MakeUpdateHandler(log, r),
			patchHandlerFunc:  patch.MakePatchHandler(log, r),
		},
	)
}

func (s *ServerAPI) Create(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
	return s.createHandlerFunc(ctx, req)
}

func (s *ServerAPI) Update(ctx context.Context, req *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	return s.updateHandlerFunc(ctx, req)
}

func (s *ServerAPI) Patch(ctx context.Context, req *proto.PatchRequest) (*proto.PatchResponse, error) {
	return s.patchHandlerFunc(ctx, req)
}

func (s *ServerAPI) Task(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	return s.taskHandlerFunc(ctx, req)
}

func (s *ServerAPI) Tasks(ctx context.Context, req *proto.TasksRequest) (*proto.TasksResponse, error) {
	return s.tasksHandlerFunc(ctx, req)
}
