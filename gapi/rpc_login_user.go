package gapi

import (
	"context"
	"fmt"

	"github.com/ymanshur/simplebank/internal/ucase"
	pb "github.com/ymanshur/simplebank/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	mtdt := server.extractMetadata(ctx)

	fmt.Println(mtdt)

	login, err := server.ucase.User.Login(ctx, ucase.LoginUserRequest{
		Username:  req.GetUsername(),
		Password:  req.Password,
		UserAgent: mtdt.UserAgent,
		ClientIp:  mtdt.ClientIP,
	})
	if err != nil {
		return nil, translationError(err)
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(&login.User),
		SessionId:             login.SessionID.String(),
		AccessToken:           login.AccessToken,
		RefreshToken:          login.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(login.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(login.RefreshTokenExpiresAt),
	}
	return rsp, nil
}
