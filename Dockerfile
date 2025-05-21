# error-handling-demo/Dockerfile

# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod ./
# If you had a go.sum (after a `go mod tidy` or `go get`)
# COPY go.sum ./
# RUN go mod download

# Copy the source code
COPY main.go ./

# Build the application
# CGO_ENABLED=0 for a static binary, GOOS=linux for cross-compilation if building on non-Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o error-demo .

# Stage 2: Create a minimal production image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/error-demo .

# Command to run the application
CMD ["./error-demo"]
