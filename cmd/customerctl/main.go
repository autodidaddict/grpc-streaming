package main

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/autodidaddict/grpc-streaming/proto"
	"google.golang.org/grpc"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("127.0.0.1:8088", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewCustomersClient(conn)

	demoTwoWayStreaming(client)
	demoServerStreaming(client)
	demoClientStreaming(client)
}

func demoClientStreaming(client pb.CustomersClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.ImportCustomers(ctx)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		cust := &pb.Customer{
			CustomerId: fmt.Sprintf("CUST-%d", i),
			GivenName:  "Random",
			Surname:    "Entropy",
		}
		if err := stream.Send(cust); err != nil {
			panic(err)
		}
		fmt.Println("Sent 1 customer for import")
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Imported all customers: %+v\n", reply)
}

func demoServerStreaming(client pb.CustomersClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.GetCustomerOrders(ctx, &pb.CustomerRequest{
		CustomerId: fmt.Sprintf("CUST-%d", 12),
	})
	if err != nil {
		panic(err)
	}
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received order on stream %+v\n", order)
	}
	fmt.Println("Downloaded all customer orders")
}

func demoTwoWayStreaming(client pb.CustomersClient) {
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
	for i := 0; i < 10; i++ {
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
