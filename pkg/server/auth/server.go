package auth

import (
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	pb "github.com/zhuangqi/biubiu/pkg/server/proto/auth"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
}

func Serve() {
	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcValidator.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcValidator.UnaryServerInterceptor()),
	)
	pb.RegisterAuthServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
