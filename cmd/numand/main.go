//Numan server executable
package main

import (
	"log"
	"net"

	"github.com/footfish/numan/api/grpc"
	"github.com/footfish/numan/internal/datastore"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	//dsn is path to sqlite db file
	dsn      = "./examples/numan-sqlite.db"
	certFile = "./examples/cert.pem"
	keyFile  = "./examples/key.pem"
	port     = ":50051"
)

func main() {

	//Database
	store := datastore.NewStore(dsn)
	defer store.Close()

	//Prep server
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	//GRPC
	log.Printf("Starting gRPC user service on %s...\n", lis.Addr().String())
	grpcServer := grpc.NewGrpcServer(creds)

	numberingServerAdapter := grpc.NewNumberingServerAdapter(store)
	userServerAdapter := grpc.NewUserServerAdapter(store)

	grpc.RegisterNumberingServer(grpcServer, numberingServerAdapter)
	grpc.RegisterUserServer(grpcServer, userServerAdapter)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
