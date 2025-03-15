# Use an official Golang runtime as a parent image
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go files and vendor directory (if exists)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o royalroadbot .

# Use a lightweight Alpine-based image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/royalroadbot /app/royalroadbot

# Expose port 8080 to the outside world
EXPOSE 8090

# Command to run the application
CMD ["./royalroadbot"]
