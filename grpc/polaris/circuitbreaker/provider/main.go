package main

import (
	"context"
	"flag"
	"github.com/tsingsun/woocoo"
	"github.com/tsingsun/woocoo-example/grpc/polaris/hellopb"
	_ "github.com/tsingsun/woocoo/contrib/polarismesh"
	"github.com/tsingsun/woocoo/rpc/grpcx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	cbEnable = flag.Bool("cb", false, "enable circuit breaker")
)

func main() {
	flag.Parse()
	app := woocoo.New()
	srv := grpcx.New(
		grpcx.WithConfiguration(app.AppConfiguration().Sub("service")),
		grpcx.WithGrpcOption(grpc.ChainUnaryInterceptor(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if *cbEnable {
					return nil, status.Error(codes.Unimplemented, "")
				}
				return handler(ctx, req)
			}),
		))
	hellopb.RegisterHelloServiceServer(srv.Engine(), &hellopb.Service{})
	app.RegisterServer(srv)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
