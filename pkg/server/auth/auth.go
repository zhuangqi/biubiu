package auth

import (
	"context"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcCtxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	pb "github.com/zhuangqi/biubiu/pkg/server/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseToken(token string) (struct{}, error) {
	return struct{}{}, nil
}

func userClaimFromToken(struct{}) string {
	return "foobar"
}

func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpcAuth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	grpcCtxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))
	newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
	return newCtx, nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{Id: 1, Username: in.Username, Token: "accesstoken_test1"}, nil
}
