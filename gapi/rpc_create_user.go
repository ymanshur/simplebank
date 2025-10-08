package gapi

import (
	"context"

	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := server.ucase.User.Create(ctx, ucase.CreateUserRequest{
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
