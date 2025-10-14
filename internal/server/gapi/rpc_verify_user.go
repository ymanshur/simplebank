package gapi

import (
	"context"

	"github.com/ymanshur/simplebank/internal/ucase"
	pb "github.com/ymanshur/simplebank/proto"
)

func (s *Server) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	user, err := s.ucase.User.Verify(ctx, ucase.VerifyUserRequest{
		EmailID:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	rsp := &pb.VerifyUserResponse{
		IsVerified: user.IsVerified,
	}
	return rsp, nil
}
