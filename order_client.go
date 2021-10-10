package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"io"
	"log"
	pb "study-grpc-client/order/order"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	searchStream, _ := c.GetOrders(ctx, &wrappers.StringValue{Value: "Google"})

	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			break
		}

		log.Printf("Search Result: ", searchOrder)
	}
}
