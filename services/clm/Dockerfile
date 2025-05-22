# Use an official Go runtime as a parent image
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# CGO_ENABLED=0 is often used for alpine builds to ensure static linking if using native Go DNS resolver.
# -ldflags="-s -w" can be used to strip debug information and reduce binary size for production.
RUN CGO_ENABLED=0 go build -o /clm-service cmd/clm/main.go

# Start a new stage from scratch for a smaller image
FROM alpine:latest
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /clm-service .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./clm-service"]
