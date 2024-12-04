# Use Go base image
FROM golang:1.23.3

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application code
COPY . .

# Copy SSL certificates (ensure they exist in your local setup)
COPY fullchain.pem /etc/ssl/certs/fullchain.pem
COPY privkey.pem /etc/ssl/private/privkey.pem

# Build the Go application
RUN go build -o server main.go

# Expose the HTTPS port
EXPOSE 8443

# Command to run the server
CMD ["./server"]
