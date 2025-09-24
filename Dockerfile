# Start from a Go base image
FROM golang:1.24.5-alpine

# Install git and mysql-client for database setup along with SSL
RUN apk add --no-cache git mysql-client openssl

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download 
#above downloads and caches project depends in go.mod file

# Copy the rest of the application source code from current Dir to /app in container
COPY . .

# Build the application binary
RUN go build -o /app/snippetbox ./cmd/web

# Expose port 4000 to the outside world
EXPOSE 4000

# Set the entrypoint script which will start the application
ENTRYPOINT ["./entrypoint.sh"]
