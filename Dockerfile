FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /app

# Install required system packages
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Create volume for persistent SQLite database
VOLUME ["/app/data"]

EXPOSE 3000

# Command to run the executable
CMD ["./main"] 