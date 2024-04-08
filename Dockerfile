# Stage 1: Build the Go binary (with .env handling)
FROM golang:1.21.3 AS builder

WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# (Optional) If using docker-compose volume mounting:
# Skip copying .env if a volume is mounted at runtime

# Otherwise (building image separately):
COPY .env ./  
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory_manager .

# Stage 2: Create a minimal container to run the Go binary
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built Go binary from the builder stage
COPY --from=builder /app/inventory_manager .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./inventory_manager"]