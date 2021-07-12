package main

import (
	"context"
	"log"
	"net"

	"mygrpc/grpc/subjects"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type subjectsServer struct {
	subjects.UnimplementedSubjectsServer
}

func (subjectsServer) Get(ctx context.Context, req *subjects.GetRequest) (*subjects.GetResponse, error) {
	id := req.GetId()
	log.Printf("'Get' request with id = %d", id)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("No metadata")
	} else {
		log.Printf("Metadata: %v", md)
	}

	rsp := subjects.GetResponse{
		User: &subjects.User{
			Id:      id + 1,
			Email:   "testuser@example.com",
			Name:    "Tester",
			Surname: "Testerson",
			Gender:  3,
		},
	}

	err := grpc.SetHeader(ctx, metadata.Pairs("X-Custom-Header-Token", "some-header-token"))
	if err != nil {
		log.Printf("Error setting header: %v", err)
	}
	err = grpc.SetTrailer(ctx, metadata.Pairs("X-Custom-Trailer-Key", "strange-trailer-value"))
	if err != nil {
		log.Printf("Error setting trailer: %v", err)
	}

	return &rsp, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	subjects.RegisterSubjectsServer(gs, &subjectsServer{})

	log.Printf("Server listening at %v", lis.Addr())
	err = gs.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
