package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func clientStreamInterceptor(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Println("=== [Client Interceptor] ", method)
	s, err := streamer(ctx, desc, cc, method, opts...)

	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

type wrappedStream struct {
	grpc.ClientStream
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

func PrintLog(m interface{}, msg string) {
	log.Printf("===== [Client Stream Interceptor] "+
		msg+" (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	PrintLog(m, "Receive a message")
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	PrintLog(m, "Send a message")
	return w.ClientStream.SendMsg(m)
}
