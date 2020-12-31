package apigateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	account "github.com/zhuangqi/biubiu/pkg/server/proto/account"
	"github.com/zhuangqi/biubiu/pkg/server/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net/http"
)

func Serve() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// grpc服务注册
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// HTTP转grpc
	err := account.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, "127.0.0.1:9101", opts)
	if err != nil {
		grpclog.Fatalf("Register handler err:%v\n", err)
	}
	err = auth.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, "127.0.0.1:9002", opts)
	if err != nil {
		grpclog.Fatalf("Register handler err:%v\n", err)
	}

	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	grpclog.Info("HTTP Listen on 8080")
	_ = http.ListenAndServe(":8081", mux)
}
