package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"time"
	pb "github.com/autodidaddict/grpc-streaming/proto"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("127.0.0.1:8088", opts ...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewCustomersClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.GetCustomerDetails(ctx)
	if err != nil {
		panic(err)
	}
	waitChannel := make(chan struct{})
	// Asynchronous stream RECEIVER
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitChannel)
				return
			}
			if err != nil {
				panic(err)
			}
			fmt.Printf("Received customer on stream %+v\n", in)
		}
	}()
	// Stream SENDER
	for i :=0; i < 10; i++ {
		customerID := fmt.Sprintf("CUST-%d", i)
		req := &pb.CustomerRequest{
			CustomerId: customerID,
		}
		if err = stream.Send(req); err != nil {
			panic(err)
		}
	}
	stream.CloseSend()
	<-waitChannel
	fmt.Println("Bi-Directional session complete.")
}