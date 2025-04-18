# Start from the official Go base image
FROM golang:1.23.1-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to download dependencies
COPY go.mod go.sum ./

# Download dependencies (Go modules)
RUN go mod download

# Copy the source code into the container
COPY . .

# Install mockgen
RUN go install go.uber.org/mock/mockgen@latest

# Generate mocks

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init --pd

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Expose port 8080 to the outside world
EXPOSE 8080

# # Command to run the executable
# CMD ["/app/main"]

# find / -maxdepth 3 -name main
# CMD ["find", "/", "-maxdepth", "3", "-name", "main"]
CMD ["/main"]