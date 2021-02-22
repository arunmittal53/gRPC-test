package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/grpc-test/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Inside client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %+v", err)
	}
	defer conn.Close()
	client := greetpb.NewGreatServiceClient(conn)

	//doUIC(client)

	//doServerStreaming(client)

	doClientStream(client)
}

func doUIC(client greetpb.GreatServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Arun",
			LastName:  "Mittal",
		},
	}
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
	fmt.Printf("Response is %v", res.Result)
}

func doServerStreaming(client greetpb.GreatServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Arun",
			LastName:  "Mittal",
		},
	}
	resStream, err := client.GreetManyTimes(context.Background(), req)
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
		fmt.Printf("response is %v\n", msg.GetResult())
	}
}

func doClientStream(client greetpb.GreatServiceClient) {
	stream, err := client.ClientGreetManyTimes(context.Background())
	if err != nil {
		log.Fatalf("error %v", err)
	}
	for i := 0; i < 5; i++ {
		req := &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Arun " + strconv.Itoa(i),
				LastName:  "Mittal",
			},
		}
		fmt.Printf("Sending %+v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Printf("Response is %s", resp)
}
