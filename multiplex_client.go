package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	ecpb "study-grpc-client/ecommerce/product"
	pb "study-grpc-client/order/order"
	"time"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderManagementClient := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	order1 := pb.Order{
		Id:          "101",
		Items:       []string{"Iphone XS", "Mac Book Pro"},
		Destination: "San Jose, CA",
		Price:       2300.00,
		Description: "test ",
	}
	res, addErr := orderManagementClient.AddOrder(ctx, &order1)

	if addErr != nil {
		got := status.Code(addErr)
		log.Printf("Error Occurred -> addOrder : %v", got)
	} else {
		log.Printf("AddOrder Response -> %v", res.Value)
	}

	ecommerceManagementClient := ecpb.NewProductInfoClient(conn)
	ecCtx, ecCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ecCancel()

	ecommerce1 := ecpb.Product{
		Id:          "101",
		Name:        "Aborcado",
		Description: "Oil tasted good",
		Price:       2.0,
	}
	ecRes, ecAddErr := ecommerceManagementClient.AddProduct(ecCtx, &ecommerce1)

	if ecAddErr != nil {
		got := status.Code(ecAddErr)
		log.Printf("Error Occurred -> addProduct : %v", got)
	} else {
		log.Printf("AddProduct Response -> %v", ecRes.Value)
	}
}
