package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/grpc-test/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet fn was inviked with %+v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetRequest, stream greetpb.GreatService_GreetManyTimesServer) error {
	firstame := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstame + strconv.Itoa(i)
		res := &greetpb.GreetResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) ClientGreetManyTimes(stream greetpb.GreatService_ClientGreetManyTimesServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &greetpb.GreetResponse{
				Result: result,
			}
			stream.SendAndClose(res)
			return nil
		}
		if err != nil {
			log.Fatalf("error %v", err)
		}
		result += req.GetGreeting().GetFirstName()
	}
}

func main() {
	fmt.Println("hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %+v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreatServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}
