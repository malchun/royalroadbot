FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o royalroadbot app/.

# ===========================
# Use a lightweight Alpine-based image for the final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/royalroadbot /app/royalroadbot
EXPOSE 8090

# Command to run the application
CMD ["./royalroadbot"]
