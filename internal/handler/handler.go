package handler

import (
	"github.com/go-kit/kit/log"
	pb "github.com/autodidaddict/grpc-streaming/proto"
	"github.com/go-kit/kit/log/level"
	"io"
)

type rpcHandler struct {
	logger log.Logger
}

func NewCustomersHandler(logger log.Logger) pb.CustomersServer {
	return &rpcHandler {
		logger: logger,
	}
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

		level.Info(r.logger).Log("msg","Received request for customer", "customer_id", customerRequest.CustomerId)

		reply := &pb.Customer{
			CustomerId: customerRequest.CustomerId,
			GivenName: "Random",
			Surname: "Entropy",
			Address: &pb.Address{
				City: "Somwewhere",
				Line1: "101 Main St",
				State: "SW",
				Zip: "99999",
			},
		}

		if err = stream.Send(reply); err != nil {
			level.Error(r.logger).Log("msg", "Failed to send reply on stream", "error", err)
			return err
		}
	}
}