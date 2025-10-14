package gapi

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ymanshur/simplebank/internal/ucase"
	pb "github.com/ymanshur/simplebank/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user *ucase.UserResponse) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func convertValidationErrors(validationErrors validation.Errors) (violations []*errdetails.BadRequest_FieldViolation) {
	for field, err := range validationErrors {
		violations = append(violations, fieldViolation(field, err))
	}
	return
}
