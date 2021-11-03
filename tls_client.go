package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	pb "study-grpc-client/ecommerce/product"
	"time"
)

var (
	url      = "localhost:50051"
	hostname = "localhost"
	crtFile  = "server.crt"
)

func GenProductParams() *pb.Product {
	ecommerce1 := pb.Product{
		Id:          "101",
		Name:        "Aborcado",
		Description: "Oil tasted good",
		Price:       2.0,
	}
	return &ecommerce1
}

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	product := GenProductParams()
	log.Printf("add a product: %v", product)
	c.AddProduct(ctx, product)
}
