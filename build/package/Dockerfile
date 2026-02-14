# Stage 1: Build
FROM golang:1.25.7-alpine3.23 AS builder

# Set the working directory
WORKDIR /app

# Install build tools
RUN apk add --no-cache make

# Install tools
RUN go install github.com/google/wire/cmd/wire@latest

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Accept version as a build argument
ARG VERSION=dev

# Generate wire files and build the application
RUN make gen
RUN CGO_ENABLED=0 GOOS=linux make build-static VERSION=${VERSION}
RUN CGO_ENABLED=0 GOOS=linux make build-worker-static VERSION=${VERSION}

# Stage 2: Runtime
FROM alpine:3.23

# Install runtime dependencies for cgo
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binaries
COPY --from=builder /app/bin/app .
COPY --from=builder /app/bin/app-worker .
# Copy the environment file
RUN mkdir -p configs
COPY --from=builder /app/configs/default.env ./configs/default.env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]
