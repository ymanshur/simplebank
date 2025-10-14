package gapi

import (
	"context"

	"github.com/ymanshur/simplebank/internal/repo"
	pb "github.com/ymanshur/simplebank/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	violations := validateVerifyUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyUserTx(ctx, repo.VerifyUserTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	rsp := &pb.VerifyUserResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyUserRequest(req *pb.VerifyUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := validateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
