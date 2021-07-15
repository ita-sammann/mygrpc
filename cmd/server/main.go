package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"mygrpc/grpc/subjects"
)

const (
	listenPort = "50051"
	listenAddr = ":" + listenPort
	localAddr  = "localhost:" + listenPort
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

func handlerFunc(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func main() {
	gs := grpc.NewServer()
	subjects.RegisterSubjectsServer(gs, &subjectsServer{})

	mux := runtime.NewServeMux()
	err := subjects.RegisterSubjectsHandlerFromEndpoint(context.Background(), mux, localAddr, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}

	log.Printf("Server started on %v", listenAddr)
	err = http.ListenAndServe(listenAddr, handlerFunc(gs, mux))
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
