package account

import (
	"github.com/zhuangqi/biubiu/pkg/constans"
	"github.com/zhuangqi/biubiu/pkg/manager"
	pb "github.com/zhuangqi/biubiu/pkg/server/proto/account"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

func Serve() {
	manager.NewGrpcServer(constans.AccountServiceHost, constans.AccountServicePort).
		Serve(func(server *grpc.Server) {
			pb.RegisterUserServiceServer(server, &Server{})
		})
}
