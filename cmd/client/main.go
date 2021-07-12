package main

import (
	"context"
	"crypto/x509"
	"log"
	"time"

	"mygrpc/grpc/subjects"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	addr = "grpc.sammann.ru:443"
	// addr = "localhost:50051"
	subjectID = 123
	authToken = "foobarbaz"
)

func getCreds() credentials.TransportCredentials {
	syspool, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}
	return credentials.NewClientTLSFromCert(syspool, "")
}

func appendAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "X-Custom-Auth-Token", authToken)
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithTransportCredentials(getCreds()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := subjects.NewSubjectsClient(conn)

	header := metadata.Pairs()
	trailer := metadata.Pairs()
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	rsp, err := c.Get(
		appendAuth(ctx),
		&subjects.GetRequest{Id: subjectID},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		log.Fatalf("Get request failed: %v", err)
	}
	log.Printf("Got response:\n\tUser: %v\n\tHeader: %v\n\tTrailer: %v", rsp.GetUser(), header, trailer)
}
