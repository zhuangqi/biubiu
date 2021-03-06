package manager

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
	"github.com/zhuangqi/biubiu/pkg/config"
	"github.com/zhuangqi/biubiu/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"strings"
	"time"
)

type checkerT func(ctx context.Context, req interface{}) error
type builderT func(ctx context.Context, req interface{}) interface{}

var (
	defaultChecker checkerT
	defaultBuilder builderT
)

type GrpcServer struct {
	ServiceName    string
	Port           int
	showErrorCause bool
	checker        checkerT
	builder        builderT
	mysqlConfig    config.MysqlConfig
}

type RegisterCallback func(*grpc.Server)

func NewGrpcServer(serviceName string, port int) *GrpcServer {
	return &GrpcServer{
		ServiceName:    serviceName,
		Port:           port,
		showErrorCause: false,
		checker:        defaultChecker,
		builder:        defaultBuilder,
	}
}

func (g *GrpcServer) ShowErrorCause(b bool) *GrpcServer {
	g.showErrorCause = b
	return g
}

func (g *GrpcServer) WithChecker(c checkerT) *GrpcServer {
	g.checker = c
	return g
}

func (g *GrpcServer) WithBuilder(b builderT) *GrpcServer {
	g.builder = b
	return g
}

func (g *GrpcServer) WithMysqlConfig(cfg config.MysqlConfig) *GrpcServer {
	g.mysqlConfig = cfg
	return g
}

func (g *GrpcServer) Serve(callback RegisterCallback, opt ...grpc.ServerOption) {
	logger.Info(nil, "Service [%s] start listen at port [%d]", g.ServiceName, g.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		err = errors.WithStack(err)
		logger.Critical(nil, "failed to listen: %+v", err)
	}
	builtinOptions := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			g.unaryServerLogInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		),
	}
	grpcServer := grpc.NewServer(append(opt, builtinOptions...)...)
	reflection.Register(grpcServer)
	callback(grpcServer)
	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Critical(nil, "%+v", err)
	}
}

var (
	jsonPbMarshaller = &jsonpb.Marshaler{
		OrigName: true,
	}
)

func (g *GrpcServer) unaryServerLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := strings.Split(info.FullMethod, "/")
		action := method[len(method)-1]
		if p, ok := req.(proto.Message); ok {
			if content, err := jsonPbMarshaller.MarshalToString(p); err != nil {
				logger.Error(ctx, "Failed to marshal proto message to string [%s]", action, err)
			} else {
				logger.Info(ctx, "Request received [%s] [%s]", action, content)
			}
		}
		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start)
		logger.Info(ctx, "Handled request [%s]  exec_time is [%s]", action, elapsed)
		return resp, err
	}
}
