# Use the official Go image to build the app
FROM golang:1.23.3

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o server main.go

# Expose the port the app runs on
EXPOSE 5000

# Command to run the application
CMD ["./server"]

