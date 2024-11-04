#!/bin/bash

# Create directories if they don't exist
mkdir -p pkg/proto/auth
mkdir -p pkg/proto/user

# Generate auth service
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/proto/auth/auth.proto

# Generate user service
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --experimental_allow_proto3_optional \
    pkg/proto/user/user.proto