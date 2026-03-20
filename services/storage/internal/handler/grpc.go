package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientpb "github.com/rado31/rabbit/proto/gen"
	"github.com/rado31/rabbit/storage/internal/model"
	"github.com/rado31/rabbit/storage/internal/service"
)

type GRPCHandler struct {
	clientpb.UnimplementedClientServiceServer
	svc *service.ClientService
}

func New(svc *service.ClientService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) CreateClient(
	ctx context.Context,
	req *clientpb.CreateClientRequest,
) (*clientpb.CreateClientResponse, error) {
	client, err := h.svc.Create(ctx, model.Client{
		Surname: req.Surname,
		Name:    req.Name,
		Age:     req.Age,
		Email:   req.Email,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &clientpb.CreateClientResponse{Client: toProto(client)}, nil
}

func (h *GRPCHandler) GetClient(
	ctx context.Context,
	req *clientpb.GetClientRequest,
) (*clientpb.GetClientResponse, error) {
	client, err := h.svc.GetByID(ctx, req.Id)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &clientpb.GetClientResponse{Client: toProto(client)}, nil
}

func toProto(c *model.Client) *clientpb.Client {
	return &clientpb.Client{
		Id:      c.ID,
		Surname: c.Surname,
		Name:    c.Name,
		Age:     c.Age,
		Email:   c.Email,
	}
}
