# Use the official Golang 1.21.4 image as a builder stage.
FROM golang:1.21.4 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
# If you have a go.sum file, uncomment the next line and make sure it's copied
# COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod file is not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o apollonion-conversation-service .

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/apollonion-conversation-service .

# Command to run the executable
CMD ["./apollonion-conversation-service"]
