package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/grpc-test/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Greet fn was inviked with %+v\n", req)
	firstNumber := req.GetNum1()
	secondNumber := req.GetNum2()
	result := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) SumServer(req *calculatorpb.SumRequest, stream calculatorpb.CalculatoeService_SumServerServer) error {
	var k int32
	k = 2
	firstNumber := req.GetNum1()
	for firstNumber > 1 {
		if firstNumber%k == 0 {
			firstNumber /= k
			res := &calculatorpb.SumResponse{
				Result: k,
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
		}
	}
	return nil
}

func (*server) AverageClient(stream calculatorpb.CalculatoeService_AverageClientServer) error {
	num := int32(0)
	sum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &calculatorpb.SumResponse{
				Result: sum / num,
			}
			stream.SendAndClose(res)
			return nil
		}
		if err != nil {
			log.Fatalf("error %+v", err)
		}
		sum += req.GetNum1()
		num++
	}
}

func main() {
	fmt.Println("hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %+v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatoeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}
