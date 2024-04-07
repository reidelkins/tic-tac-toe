# Start from a base image with Go installed
FROM golang:1.22 as builder

# Install Air for live reloading
RUN go install github.com/cosmtrek/air@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Copy the Air configuration file
COPY .air.toml ./

# Run the Air daemon
CMD ["air"]




# # Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cmd/server/main ./cmd/server


# # Start a new stage from scratch
# FROM alpine:latest  

# # Set the Current Working Directory inside the container
# WORKDIR /root/

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/cmd/server/main .

# # Expose port 8080 to the outside world
# EXPOSE 8080

# # Command to run the executable
# CMD ["./main"]
