package grpc

import (
	context "context"
	"errors"
	"fmt"
	"log"

	"github.com/footfish/numan"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// NewGrpcServer creates a new grpc.Server
func NewGrpcServer(creds credentials.TransportCredentials) *grpc.Server {
	return grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(authServerInterceptor))
}

// NewGrpcClient creates a new grpc client connection
func NewGrpcClient(ctx context.Context, address string, creds credentials.TransportCredentials) *grpc.ClientConn {
	// Set up a connection to the server.
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(authClientInterceptor), grpc.WithBlock())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatalf("gRPC connect timeout (check server is running)")
		}
		log.Fatalf("gRPC connect error: %s", err)

	}
	return conn
}

//authServerInterceptor copies a token from gRPC metadata to context
func authServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok && len(meta[numan.AuthTokenField]) == 1 {
		ctx = context.WithValue(ctx, numan.AuthTokenField, meta[numan.AuthTokenField][0])
	}
	return handler(ctx, req)
}

//authClientInterceptor copies a token from context to gRPC metadata
func authClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	// Intercept logic
	if token := ctx.Value(numan.AuthTokenField); token != nil { //add auth token to RPC metadata
		ctx = metadata.AppendToOutgoingContext(ctx, numan.AuthTokenField, fmt.Sprintf("%v", token))
	}
	// Calls the invoker to execute RPC
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
