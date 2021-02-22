package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpc-test/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Inside client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %+v", err)
	}
	defer conn.Close()
	client := calculatorpb.NewCalculatoeServiceClient(conn)

	//doUniary(client)

	//doServerStream(client)

	doClientStream(client)
}

func doUniary(client calculatorpb.CalculatoeServiceClient) {
	req := &calculatorpb.SumRequest{
		Num1: 10,
		Num2: 30,
	}
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
	fmt.Printf("Response is %v\n", res.Result)
}

func doServerStream(client calculatorpb.CalculatoeServiceClient) {

	req := &calculatorpb.SumRequest{
		Num1: 100,
		Num2: 30,
	}
	resStream, err := client.SumServer(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to serve: %+v", err)
		}
		fmt.Printf("%d\t", msg.GetResult())
	}
}

func doClientStream(client calculatorpb.CalculatoeServiceClient) {
	stream, err := client.AverageClient(context.Background())
	if err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
	for i := 0; i < 5; i++ {
		req := &calculatorpb.SumRequest{
			Num1: int32(i),
		}
		stream.Send(req)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Printf("Response is %s", resp)
}
