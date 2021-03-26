//Numan server executable
package main

import (
	"log"
	"net"

	"github.com/footfish/numan/api/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	//dsn is path to sqlite db file
	dsn      = "./examples/numan-sqlite.db"
	certFile = "./examples/cert.pem"
	keyFile  = "./examples/key.pem"
	port     = ":50051"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	log.Printf("Starting gRPC user service on %s...\n", lis.Addr().String())

	grpcServer, numanServerAdapter := grpc.NewGrpcServer(dsn, creds)
	grpc.RegisterNumanServer(grpcServer, numanServerAdapter)
	defer grpc.CloseServerAdapter(numanServerAdapter)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
