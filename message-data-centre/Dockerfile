# Use the official Golang image to create a build artifact.
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Install protoc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.19.0/protoc-3.19.0-linux-x86_64.zip && \
    unzip protoc-3.19.0-linux-x86_64.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm -rf protoc-3.19.0-linux-x86_64.zip protoc3

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set Go environment variables
ENV GO111MODULE=on \
    GOFLAGS=-mod=vendor

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Compile protobuf files
RUN protoc --go_out=. --go-grpc_out=. proto/service.proto

# Ensure consistency in vendoring
RUN go mod tidy

RUN go mod vendor

# Build the Go app
RUN go build -o main ./server

# Start a new stage with a compatible GLIBC version
FROM debian:stable-slim

# Install necessary dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Command to run the executable
ENTRYPOINT ["/app/main"]
