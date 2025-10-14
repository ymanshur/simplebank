package gapi

import (
	"context"

	"github.com/ymanshur/simplebank/internal/ucase"
	pb "github.com/ymanshur/simplebank/proto"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := s.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	user, err := s.ucase.User.Update(ctx, ucase.UpdateUserRequest{
		Auth: ucase.AuthRequest{
			Username: authPayload.Username,
			Role:     authPayload.Role,
		},
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
