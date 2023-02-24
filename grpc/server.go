package main

import (
	"context"
	"fmt"
	"github.com/tsingsun/woocoo"
	"github.com/tsingsun/woocoo-example/funcs"
	"github.com/tsingsun/woocoo-example/grpc/api"
	"github.com/tsingsun/woocoo/contrib/telemetry"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/rpc/grpcx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"time"

	_ "github.com/tsingsun/woocoo/contrib/telemetry/otelgrpc"
	_ "github.com/tsingsun/woocoo/rpc/grpcx/registry/etcd3"
)

var tracer = otel.Tracer(conf.Global().AppName())

func main() {
	app := woocoo.New()
	otelCnf := app.AppConfiguration().Sub("otel")
	otelcfg := telemetry.NewConfig(otelCnf,
		telemetry.WithTracerProviderOptions(funcs.ZipkinTracer(otelCnf)...),
		telemetry.WithMeterProviderOptions(funcs.PrometheusProvider(otelCnf)...),
	)
	defer otelcfg.Shutdown()

	srv := grpcx.New(
		grpcx.WithConfiguration(app.AppConfiguration().Sub("grpc")),
		grpcx.WithGrpcLogger(),
	)
	api.RegisterHelloServiceServer(srv.Engine(), &server{})
	app.RegisterServer(srv)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// server is used to implement api.HelloServiceServer.
type server struct {
	api.HelloServiceServer
}

// SayHello implements api.HelloServiceServer.
func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	log.Printf("Received: %v\n", in.GetGreeting())
	s.workHard(ctx)
	time.Sleep(50 * time.Millisecond)

	return &api.HelloResponse{Reply: "Hello " + in.Greeting}, nil
}

func (s *server) workHard(ctx context.Context) {
	_, span := tracer.Start(ctx, "workHard",
		trace.WithAttributes(attribute.String("extra.key", "extra.value")))
	defer span.End()

	time.Sleep(50 * time.Millisecond)
}

func (s *server) SayHelloServerStream(in *api.HelloRequest, out api.HelloService_SayHelloServerStreamServer) error {
	log.Printf("Received: %v\n", in.GetGreeting())

	for i := 0; i < 5; i++ {
		err := out.Send(&api.HelloResponse{Reply: "Hello " + in.Greeting})
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(i*50) * time.Millisecond)
	}

	return nil
}

func (s *server) SayHelloClientStream(stream api.HelloService_SayHelloClientStreamServer) error {
	i := 0

	for {
		in, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Non EOF error: %v\n", err)
			return err
		}

		log.Printf("Received: %v\n", in.GetGreeting())
		i++
	}

	time.Sleep(50 * time.Millisecond)

	return stream.SendAndClose(&api.HelloResponse{Reply: fmt.Sprintf("Hello (%v times)", i)})
}

func (s *server) SayHelloBidiStream(stream api.HelloService_SayHelloBidiStreamServer) error {
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Non EOF error: %v\n", err)
			return err
		}

		time.Sleep(50 * time.Millisecond)

		log.Printf("Received: %v\n", in.GetGreeting())
		err = stream.Send(&api.HelloResponse{Reply: "Hello " + in.Greeting})

		if err != nil {
			return err
		}
	}

	return nil
}
