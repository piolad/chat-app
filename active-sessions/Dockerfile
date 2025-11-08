# syntax=docker/dockerfile:1.6

############################
# 1) Builder image
############################
ARG RUST_VERSION=1.85
FROM rust:${RUST_VERSION}-slim AS builder

# System deps only needed to *build* (not at runtime)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      pkg-config \
      libssl-dev \
      protobuf-compiler \
      ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# ---- Cache dependencies: copy manifests first
COPY Cargo.toml Cargo.lock ./
# If using a workspace, copy the workspace manifest and each member's Cargo.toml
# COPY Cargo.toml Cargo.lock ./
# COPY crates/some_crate/Cargo.toml crates/some_crate/Cargo.toml
# ... (repeat for members)

# Create a dummy src to satisfy cargo build and cache deps
RUN mkdir -p src && printf 'fn main(){println!("placeholder");}\n' > src/main.rs

# Build once to cache dependencies (uses BuildKit cache mounts)
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/usr/local/cargo/git \
    cargo build --release

# ---- Now copy real sources and build your app
# (Replace this with a more selective COPY if you like)
COPY . .

# Name your binary here (must match `[[bin]] name` or package name)
ARG BIN=active-sessions

# Build the real thing (again using cache mounts)
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/usr/local/cargo/git \
    cargo build --release --bin ${BIN}

############################
# 2) Runtime image
############################
FROM debian:bookworm-slim AS runtime

# Install only what the binary needs at runtime.
# Many apps only need ca-certificates; add libssl3, tzdata, etc. if your binary links to them.
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy just the compiled binary
ARG BIN=active-sessions
COPY --from=builder /app/target/release/${BIN} /usr/local/bin/${BIN}

# Drop privileges
RUN useradd -u 10001 -m nonroot
USER 10001

ENV RUST_LOG=info
CMD ["/usr/local/bin/active-sessions"]
