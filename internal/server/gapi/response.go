package gapi

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/ymanshur/simplebank/internal/typex"
	"github.com/ymanshur/simplebank/internal/ucase"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}

func translationError(err error) error {
	errCause := errors.Cause(err)
	switch errCause {
	case ucase.ErrUserUniqueViolation:
		return status.Error(codes.AlreadyExists, errCause.Error())
	}

	switch errCause := errCause.(type) {
	case validation.Errors:
		return invalidArgumentError(convertValidationErrors(errCause))
	case typex.ErrDataNotFound:
		return status.Error(codes.NotFound, errCause.Error())
	case typex.ErrUnAuthorized:
		return status.Error(codes.Unauthenticated, errCause.Error())
	case typex.ErrForbidden:
		return status.Error(codes.PermissionDenied, errCause.Error())
	default:
		return status.Error(codes.Internal, errCause.Error())
	}
}
