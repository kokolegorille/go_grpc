package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"go-grpc/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting client.")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
	doServerStreaming(c)
	doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Koko",
			LastName:  "Le Gorille",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v\n", err)
	}
	log.Printf("Response from Greet: %v\n", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Koko",
			LastName:  "Le Gorille",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v\n", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Koko",
				LastName:  "Le Gorille",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mamita",
				LastName:  "Amoi",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Meli7",
				LastName:  "Ponchou",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lito",
				LastName:  "Sandal",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jade",
				LastName:  "Pervenche",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v\n", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi streaming RPC...")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone: %v\n", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Koko",
				LastName:  "Le Gorille",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mamita",
				LastName:  "Amoi",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Meli7",
				LastName:  "Ponchou",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lito",
				LastName:  "Sandal",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jade",
				LastName:  "Pervenche",
			},
		},
	}

	waitc := make(chan struct{})

	// Sends a bunch of messages
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
		}
		stream.CloseSend()
	}()

	// Receives a bunch of messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
			}
			fmt.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
