package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "study-grpc-client/ecommerce/product"
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
	c := pb.NewProductInfoClient(conn)
	params := &pb.Product{Name: "Apple iPhone 12 Pro", Description: "Meet Apple iPhone 12 Pro", Price: float32(1000.0)}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AddProduct(ctx, params)
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID %s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", product)
	}
	log.Printf("Product: ", product.String())
}
