#! /bin/bash 
# Builds gRPC go code from protobuf file 
# Requires protoc - https://grpc.io/docs/protoc-installation/
cd ../api/grpc
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  numan.proto