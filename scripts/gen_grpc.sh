#! /bin/bash 
# Builds gRPC go code from protobuf file 
# Requires protoc - https://grpc.io/docs/protoc-installation/
# install protoc binary as above then run the following;
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

cd ../api/grpc
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  *.proto