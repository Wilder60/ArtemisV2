package adpater

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wilder60/KeyRing/security"

	"github.com/Wilder60/KeyRing/internal/grpc/service"
	"google.golang.org/grpc/metadata"
)

type ServiceGrpcImpl struct {
}

func NewServiceGrpcImpl() *ServiceGrpcImpl {
	return &ServiceGrpcImpl{}
}

// DeleteUser will go through the database and delete all event assoicated with that user
// The message that is sent from the account service will contain the id related to the user
func (serviceImpl *ServiceGrpcImpl) DeleteUser(ctx context.Context,
	in *service.DeleteUserRequest) (*service.UserResponse, error) {
	token, err := extractAuthorizationToken(ctx)
	if err != nil {
		return &service.UserResponse{
			Error: err.Error(),
		}, err
	}

	if err = security.Validate(token); err != nil {
		return &service.UserResponse{
			Error: err.Error(),
		}, err
	}

	fmt.Printf("You've deleted User with Id %d\n", in.Id)
	return &service.UserResponse{
		Error: "",
	}, nil
}

func extractAuthorizationToken(ctx context.Context) (string, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	var errorMessage string
	if !ok {
		errorMessage = "No Metadata present in request"
		return "", errors.New(errorMessage)
	}
	token, ok := headers["authorization"]
	return token[0], nil
}
