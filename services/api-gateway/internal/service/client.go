package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rado31/rabbit/api-gateway/internal/model"
	clientpb "github.com/rado31/rabbit/proto/gen"
)

type ClientService struct {
	grpc clientpb.ClientServiceClient
}

func New(grpc clientpb.ClientServiceClient) *ClientService {
	return &ClientService{grpc: grpc}
}

func (s *ClientService) Create(
	ctx context.Context,
	req model.CreateClientRequest,
) (*model.Client, error) {
	resp, err := s.grpc.CreateClient(ctx, &clientpb.CreateClientRequest{
		Surname: req.Surname,
		Name:    req.Name,
		Age:     req.Age,
		Email:   req.Email,
	})

	if err != nil {
		return nil, fmt.Errorf("storage.CreateClient: %w", err)
	}

	return toModel(resp.Client), nil
}

func (s *ClientService) GetByID(ctx context.Context, id int32) (*model.Client, error) {
	resp, err := s.grpc.GetClient(ctx, &clientpb.GetClientRequest{Id: id})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("client %d: %w", id, model.ErrNotFound)
		}

		return nil, fmt.Errorf("storage.GetClient: %w", err)
	}

	return toModel(resp.Client), nil
}

func toModel(c *clientpb.Client) *model.Client {
	return &model.Client{
		ID:      c.Id,
		Surname: c.Surname,
		Name:    c.Name,
		Age:     c.Age,
		Email:   c.Email,
	}
}
