# First stage: build the Rust application
FROM rust:latest as builder

# Install the protobuf compiler - needed for static builds
RUN apt-get update && apt-get install -y protobuf-compiler

# Add the x86_64-unknown-linux-musl target - for debian
RUN rustup target add x86_64-unknown-linux-musl

RUN mkdir -p /home/app
WORKDIR /home/app
COPY . .

RUN cargo build --release --target x86_64-unknown-linux-musl

# Second stage: create a minimal image with the compiled binary
FROM debian:buster-slim
RUN mkdir -p /home/app
WORKDIR /home/app

# Correct path to the binary
COPY --from=builder /home/app/target/x86_64-unknown-linux-musl/release/auth-service .
CMD ["./auth-service"]
