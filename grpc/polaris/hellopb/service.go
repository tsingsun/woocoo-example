package hellopb

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
)

// Service is used to implement api.HelloServiceServer.
type Service struct {
	HelloServiceServer
}

// SayHello implements api.HelloServiceServer.
func (s *Service) SayHello(ctx context.Context, in *HelloRequest) (*HelloResponse, error) {
	log.Printf("Received: %v\n", in.GetGreeting())
	s.workHard(ctx)
	time.Sleep(50 * time.Millisecond)

	return &HelloResponse{Reply: "Hello " + in.Greeting}, nil
}

func (s *Service) workHard(ctx context.Context) {
	time.Sleep(50 * time.Millisecond)
}

func (s *Service) SayHelloServerStream(in *HelloRequest, out HelloService_SayHelloServerStreamServer) error {
	log.Printf("Received: %v\n", in.GetGreeting())

	for i := 0; i < 5; i++ {
		err := out.Send(&HelloResponse{Reply: "Hello " + in.Greeting})
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(i*50) * time.Millisecond)
	}

	return nil
}

func (s *Service) SayHelloClientStream(stream HelloService_SayHelloClientStreamServer) error {
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

	return stream.SendAndClose(&HelloResponse{Reply: fmt.Sprintf("Hello (%v times)", i)})
}

func (s *Service) SayHelloBidiStream(stream HelloService_SayHelloBidiStreamServer) error {
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
		err = stream.Send(&HelloResponse{Reply: "Hello " + in.Greeting})

		if err != nil {
			return err
		}
	}

	return nil
}
