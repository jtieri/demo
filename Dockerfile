# Use the official Go image to build the binary
FROM golang:1.23 AS builder

# Set environment variables for Go modules and project path
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and download dependencies (caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go binary
WORKDIR /app/cmd/demod
RUN go build -o /demod

# Create a minimal image for running the Go binary
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /demod bin/demod

WORKDIR /root

# rest server, tendermint p2p, tendermint rpc
# EXPOSE 1317 26656 26657

# Set the entrypoint
ENTRYPOINT ["/bin/demod"]
