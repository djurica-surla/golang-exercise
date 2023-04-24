# Use the official Go base image
FROM golang:1.19-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code to the container
COPY . .

# Build the Go application
RUN go build -o /go/bin/app ./cmd/server

# Set the command to run the Go application
CMD ["/go/bin/app"]