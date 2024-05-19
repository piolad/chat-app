FROM golang:latest

# Install the Protocol Buffers compiler
RUN apt-get update && apt-get install -y protobuf-compiler

# Install the protoc-gen-go plugin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod .
COPY go.sum .

# Install dependencies
RUN go mod download

# Copy the main.go file
COPY main.go .

# Copy the .proto file
COPY protos/*.proto ./protos/

# Generate Go code from the .proto file
RUN protoc --go_out=. --go_opt=paths=source_relative ./protos/messagecentre.proto

# Build the Go application
RUN go build -o bin/main .

# Set the command to run when the container starts
CMD ["/app/bin/main"]
