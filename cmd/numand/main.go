//Numan server executable
package main

import (
	"fmt"
	"log"
	"net"

	"github.com/footfish/numan/api/grpc"
	"github.com/footfish/numan/internal/datastore"
	_ "github.com/joho/godotenv/autoload" //autoloads .env file
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	var conf struct {
		Dsn     string
		Port    int `envconfig:"default=50051,optional"`
		TlsCert string
		TlsKey  string
	}

	//Load conf from environmental vars (.env file autoloaded if present)
	if err := envconfig.Init(&conf); err != nil {
		log.Fatalf("Failed to load required environmental variables for config: %v", err)
	}

	//Database
	store := datastore.NewStore(conf.Dsn)
	defer store.Close()

	//Prep server
	creds, err := credentials.NewServerTLSFromFile(conf.TlsCert, conf.TlsKey)
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
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
		log.Fatalf("gRPC server failed to serve: %v", err)
	}

}
