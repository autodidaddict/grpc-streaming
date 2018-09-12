package handler

import (
	"fmt"
	"io"
	"time"

	pb "github.com/autodidaddict/grpc-streaming/proto"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type rpcHandler struct {
	logger log.Logger
}

// NewCustomersHandler creates an instance of a gRPC endpoint handler for the customers service
func NewCustomersHandler(logger log.Logger) pb.CustomersServer {
	return &rpcHandler{
		logger: logger,
	}
}

func (r *rpcHandler) ImportCustomers(stream pb.Customers_ImportCustomersServer) error {
	var recordCount int32
	startTime := time.Now()

	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.ImportSummary{
				TotalCustomers: recordCount,
				ErrorCount:     95781,
				ElapsedTime:    int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		level.Info(r.logger).Log("msg", "Received customer import", "customer", customer)
		recordCount++
		time.Sleep(500 * time.Millisecond)
	}
}

func (r *rpcHandler) GetCustomerOrders(req *pb.CustomerRequest, stream pb.Customers_GetCustomerOrdersServer) error {
	for i := 0; i < 10; i++ {
		order := &pb.Order{
			OrderId: fmt.Sprintf("ORDER-%d", i),
			Qty:     int32(i + 5),
		}
		if err := stream.Send(order); err != nil {
			level.Error(r.logger).Log("msg", "Failed to send reply on stream", "error", err)
			return err
		}
	}
	return nil
}

func (r *rpcHandler) GetCustomerDetails(stream pb.Customers_GetCustomerDetailsServer) error {
	level.Info(r.logger).Log("msg", "Beginning bi-directional streaming session.")

	for {
		customerRequest, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			level.Error(r.logger).Log("msg", "Failed to read from stream", "error", err)
			return err
		}

		level.Info(r.logger).Log("msg", "Received request for customer", "customer_id", customerRequest.CustomerId)

		reply := &pb.Customer{
			CustomerId: customerRequest.CustomerId,
			GivenName:  "Random",
			Surname:    "Entropy",
			Address: &pb.Address{
				City:  "Somwewhere",
				Line1: "101 Main St",
				State: "SW",
				Zip:   "99999",
			},
		}

		if err = stream.Send(reply); err != nil {
			level.Error(r.logger).Log("msg", "Failed to send reply on stream", "error", err)
			return err
		}
	}
}
