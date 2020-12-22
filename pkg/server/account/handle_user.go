package account

import (
	"context"
	pb "github.com/zhuangqi/biubiu/pkg/server/proto/account"
)

func (s *Server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	return &pb.CreateUserReply{Id: 1, Username: in.Username}, nil
}
