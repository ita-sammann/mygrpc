package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"mygrpc/grpc/subjects"
)

const (
	listenAddr = ":50051"
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

func serveGRPC(listener net.Listener) {
	gs := grpc.NewServer()
	subjects.RegisterSubjectsServer(gs, &subjectsServer{})

	go func() {
		log.Printf("gRPC server listening at %v", listener.Addr())
		err := gs.Serve(listener)
		if err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
}

func serveHTTP(listener net.Listener) {
	mux := runtime.NewServeMux()
	err := subjects.RegisterSubjectsHandlerFromEndpoint(context.Background(), mux, listenAddr, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}

	go func() {
		log.Printf("HTTP server listening at %v", listener.Addr())
		err = http.Serve(listener, mux)
		if err != nil {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()
}

func main() {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Setup multiplexer for gRPC and HTTP servers
	mux := cmux.New(lis)

	lisGrpc := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	lisHttp := mux.Match(cmux.Any())

	serveGRPC(lisGrpc)
	serveHTTP(lisHttp)

	log.Printf("Multiplexer started on %v", lis.Addr())
	err = mux.Serve()
	if err != nil {
		log.Fatalf("Failed to serve multiplexer: %v", err)
	}
}
