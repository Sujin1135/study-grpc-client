package main

import (
	"io"
	"log"
	pb "study-grpc-client/order/order"
)

const (
	address = "localhost:50051"
)

func asyncClientBidirectionalRPC(streamProcOrder pb.OrderManagement_ProcessOrdersClient,
	c chan struct{}) {
	for {
		combinedShipment, errProcOrder := streamProcOrder.Recv()
		if errProcOrder == io.EOF {
			break
		}
		if errProcOrder != nil {
			log.Fatalf("error is ", errProcOrder)
		}
		log.Printf("Id", combinedShipment.Id)
		log.Printf("Status", combinedShipment.Status)
		log.Printf("Order list ", combinedShipment.OrderList)
	}
	<-c
}

//func main() {
//	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithStreamInterceptor(clientStreamInterceptor))
//	if err != nil {
//		log.Fatalf("did not connect: %v", err)
//	}
//	defer conn.Close()
//
//	c := pb.NewOrderManagementClient(conn)
//	deadline := time.Now().Add(time.Duration(2 * time.Second))
//	ctx, cancel := context.WithDeadline(context.Background(), deadline)
//	defer cancel()
//
//	// search example as below
//	//searchStream, _ := c.GetOrders(ctx, &wrappers.StringValue{Value: "Google"})
//	//
//	//for {
//	//	searchOrder, err := searchStream.Recv()
//	//	if err == io.EOF {
//	//		break
//	//	}
//	//
//	//	log.Printf("Search Result: ", searchOrder)
//	//}
//
//	// deadline example
//	// start an update example as below
//	updateStream, err := c.UpdateOrders(ctx)
//
//	if err != nil {
//		log.Fatalf("%v.UpdateOrders(_) = _, %v", c, err)
//	}
//
//	updOrder1 := pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Google Pixel Book"}, Destination: "Mountain View, CA", Price: 1100.00}
//	updOrderError := pb.Order{Id: "-1", Items: []string{"Error", "Google Pixel Book"}, Destination: "Mountain View, CA", Price: 1100.00}
//	updOrder2 := pb.Order{Id: "103", Items: []string{"Apple Watch S4", "Mac Book Pro", "iPad Pro"}, Destination: "San Jose, CA", Price: 2800.00}
//	updOrder3 := pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub", "iPad Mini"}, Destination: "Mountain View, CA", Price: 2200.00}
//
//	if err := updateStream.Send(&updOrder1); err != nil {
//		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder1, err)
//	}
//
//	if err := updateStream.Send(&updOrderError); err != nil {
//		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder1, err)
//	}
//
//	if err := updateStream.Send(&updOrder2); err != nil {
//		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder2, err)
//	}
//
//	if err := updateStream.Send(&updOrder3); err != nil {
//		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder3, err)
//	}
//
//	updateRes, err := updateStream.CloseAndRecv()
//	if err != nil {
//		errorCode := status.Code(err)
//		log.Printf("Invalid Argument Error : %s", errorCode)
//		errorStatus := status.Convert(err)
//		errDetails := errorStatus.Details()
//
//		log.Printf("Error status : %s", errorStatus)
//
//		for _, d := range errDetails {
//			log.Printf("In for loop")
//			switch info := d.(type) {
//			case *epb.BadRequest_FieldViolation:
//				log.Printf("Request Field Invalid: %s", info)
//			default:
//				log.Printf("Unexcepted error type : %s", info)
//			}
//		}
//		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateStream, err, nil)
//	}
//	log.Printf("Update Orders Res : %s", updateRes)
//
//	// process order example as below
//	//streamProcOrder, _ := c.ProcessOrders(ctx)
//	//if err := streamProcOrder.Send(&wrappers.StringValue{Value: "102"}); err != nil {
//	//	log.Fatalf("%v.Send(%v) = %v", c, "102", err)
//	//}
//	//
//	//if err := streamProcOrder.Send(&wrappers.StringValue{Value: "103"}); err != nil {
//	//	log.Fatalf("%v.Send(%v) = %v", c, "103", err)
//	//}
//	//
//	//if err := streamProcOrder.Send(&wrappers.StringValue{Value: "104"}); err != nil {
//	//	log.Fatalf("%v.Send(%v) = %v", c, "104", err)
//	//}
//	//
//	//channel := make(chan struct{})
//	//go asyncClientBidirectionalRPC(streamProcOrder, channel)
//	//time.Sleep(time.Millisecond * 1000)
//	//
//	//if err := streamProcOrder.Send(&wrappers.StringValue{Value: "101"}); err != nil {
//	//	log.Fatalf("%v.Send(%v) = %v", c, "101", err)
//	//}
//	//
//	//if err := streamProcOrder.CloseSend(); err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//<-channel
//}
