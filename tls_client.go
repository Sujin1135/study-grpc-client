package main

import (
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	pb "study-grpc-client/ecommerce/product"
	"time"
)

const (
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
	auth := oauth.NewOauthAccess(fetchToken())

	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %s", err)
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(auth),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	product := GenProductParams()
	log.Printf("add a product: %v", product)
	id, err := c.AddProduct(ctx, product)

	if err != nil {
		log.Fatalf("failed to add a product by id %s, %s", product.Id, err)
	} else {
		log.Printf("added id is %s", id)
	}
}

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
