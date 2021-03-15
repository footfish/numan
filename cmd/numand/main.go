//Numan server executable
package main

import (
	"log"
	"net"

	"github.com/footfish/numan/api/grpc"
)

const (
	//dsn is path to sqlite db file
	dsn  = "./examples/numan-sqlite.db"
	port = ":50051"
)

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	log.Printf("Starting gRPC user service on %s...\n", lis.Addr().String())

	grpcServer, numanServerAdapter := grpc.NewGrpcServer(dsn)
	grpc.RegisterNumanServer(grpcServer, numanServerAdapter)
	defer grpc.CloseServerAdapter(numanServerAdapter)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
