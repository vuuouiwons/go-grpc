package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "go-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedSimpleServiceServer
}

func (s *server) Normal(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	fmt.Printf("[Server] Normal received: %s\n", req.GetMessage())

	if req.GetError() {
		return nil, status.Errorf(codes.InvalidArgument, "server received error param=true for message: %s", req.GetMessage())
	}

	return &pb.Response{
		Message: "processed: " + req.GetMessage(),
	}, nil
}

func (s *server) ServerStream(req *pb.Request, stream pb.SimpleService_ServerStreamServer) error {
	fmt.Printf("[Server] ServerStream received: %s\n", req.GetMessage())

	if req.GetError() {
		return status.Errorf(codes.InvalidArgument, "server received error param=true for message: %s", req.GetMessage())
	}

	for i := 1; i <= 3; i++ {
		time.Sleep(500 * time.Millisecond)
		response := &pb.Response{
			Message: fmt.Sprintf("[%d] stream: %s", i, req.GetMessage()),
		}
		fmt.Printf("[Server] ServerStream sending: %s\n", response.GetMessage())
		if err := stream.Send(response); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ClientStream(stream pb.SimpleService_ClientStreamServer) error {
	var count int

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Response{
				Message: fmt.Sprintf("received total %d messages", count),
			})
		}
		if err != nil {
			return err
		}

		fmt.Printf("[Server] ClientStream received: %s\n", req.GetMessage())

		if req.GetError() {
			return status.Errorf(codes.InvalidArgument, "server received error param=true on message: %s", req.GetMessage())
		}

		count++
	}
}

func (s *server) BiDirectionalStream(stream pb.SimpleService_BiDirectionalStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("[Server] BiDirectionalStream received: %s\n", req.GetMessage())

		if req.GetError() {
			return status.Errorf(codes.InvalidArgument, "server received error param=true on message: %s", req.GetMessage())
		}

		time.Sleep(500 * time.Millisecond)
		err = stream.Send(&pb.Response{
			Message: "echo: " + req.GetMessage(),
		})
		if err != nil {
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleServiceServer(grpcServer, &server{})

	log.Println("server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
