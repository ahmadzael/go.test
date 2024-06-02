# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS build

# Set environment variables for Go
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifest and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Stage 2: Create the production image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go binary from the build stage
COPY --from=build /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the Go binary
CMD ["./main"]
