package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	echopb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	ecpb "study-grpc-client/ecommerce/product"
	pb "study-grpc-client/order/order"
	"time"
)

func MakeOrderParams() pb.Order {
	order := pb.Order{
		Id:          "101",
		Items:       []string{"Iphone XS", "Mac Book Pro"},
		Destination: "San Jose, CA",
		Price:       2300.00,
		Description: "test ",
	}
	return order
}

func AddOrder(conn *grpc.ClientConn, params pb.Order) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"kn", "vn",
	)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)
	ctxA := metadata.AppendToOutgoingContext(mdCtx, "k1", "v1", "k1", "v2", "k2", "v3")
	orderManagementClient := pb.NewOrderManagementClient(conn)
	res, addErr := orderManagementClient.AddOrder(ctxA, &params)

	if addErr != nil {
		got := status.Code(addErr)
		log.Printf("Error Occurred -> addOrder : %v", got)
	} else {
		log.Printf("AddOrder Response -> %v", res.Value)
	}
}

func MakeProductParams() ecpb.Product {
	ecommerce1 := ecpb.Product{
		Id:          "101",
		Name:        "Aborcado",
		Description: "Oil tasted good",
		Price:       2.0,
	}
	return ecommerce1
}

func AddProduct(conn *grpc.ClientConn, params ecpb.Product) {
	ecommerceManagementClient := ecpb.NewProductInfoClient(conn)
	ecCtx, ecCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ecCancel()

	var header, trailer metadata.MD
	ecRes, ecAddErr := ecommerceManagementClient.AddProduct(
		ecCtx,
		&params,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	log.Printf("header is %v", header)
	log.Printf("trailer is %v", trailer)

	if ecAddErr != nil {
		got := status.Code(ecAddErr)
		log.Printf("Error Occurred -> addProduct : %v", got)
	} else {
		log.Printf("AddProduct Response -> %v", ecRes.Value)
	}
}

func MakeUpdateOrderParams() []pb.Order {
	orders := []pb.Order{
		pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Google Pixel Book"}, Destination: "Mountain View, CA", Price: 1100.00},
		pb.Order{Id: "103", Items: []string{"Apple Watch S4", "Mac Book Pro", "iPad Pro"}, Destination: "San Jose, CA", Price: 2800.00},
	}
	return orders
}

func UpdateOrder(conn *grpc.ClientConn, params []pb.Order) {
	c := pb.NewOrderManagementClient(conn)
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
	)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)
	ctxA := metadata.AppendToOutgoingContext(mdCtx, "k1", "v1", "k1", "v2", "k2", "v3")
	stream, err := c.UpdateOrders(ctxA)

	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v", c, err)
	}

	for _, value := range params {
		if err := stream.Send(&value); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, value, err)
		}
	}

	updateRes, err := stream.CloseAndRecv()
	if err != nil {
		errorCode := status.Code(err)
		log.Printf("Invalid Argument Error : %s", errorCode)
		errorStatus := status.Convert(err)

		log.Printf("Error status : %s", errorStatus)
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Update Orders Res : %s", updateRes)

	header, err := stream.Header()
	trailer := stream.Trailer()

	log.Printf("stream header is %v", header)
	log.Printf("stream trailer is %v", trailer)
}

func callUnaryEcho(c echopb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &echopb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

// TODO: round robin case
func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := echopb.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/load_balancing")
	}
}

//func main() {
//	conn, err := grpc.Dial(
//		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
//		grpc.WithBalancerName("pick_first"),
//		grpc.WithInsecure(),
//	)
//	if err != nil {
//		log.Fatalf("did not connect: %v", err)
//	}
//	defer conn.Close()
//
//	log.Println("==== Calling AddProduct/UpdateOrder with pick_first ====")
//
//	makeRPCs(conn, 10)
//
//	// TODO: round robin
//	// round_robin 정책으로 다른 ClientConn을 만든다.
//	roundrobinConn, err := grpc.Dial(
//		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
//		// "example:///lb.example.grpc.io"
//		grpc.WithBalancerName("round_robin"),
//		grpc.WithInsecure(),
//	)
//	if err != nil {
//		log.Fatalf("did not connect: %v", err)
//	}
//	defer roundrobinConn.Close()
//
//	log.Println("==== Calling helloworld.Greeter/SayHello " +
//		"with round_robin ====")
//	makeRPCs(roundrobinConn, 10)
//
//	// AddOrder(conn, MakeOrderParams())
//	//AddProduct(conn, MakeProductParams())
//	//UpdateOrder(conn, MakeUpdateOrderParams())
//}
