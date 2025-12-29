# Stage 1: Build
FROM golang:1.25.5-alpine3.23 AS builder

# Set the working directory
WORKDIR /app

# Install build tools
RUN apk add --no-cache make gcc musl-dev

# Install migrate tool
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
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
RUN CGO_ENABLED=1 GOOS=linux make build-static VERSION=${VERSION}
RUN CGO_ENABLED=1 GOOS=linux make build-worker-static VERSION=${VERSION}

# Run the migrations
RUN make migrate-up

# Stage 2: Runtime
FROM alpine:3.23

# Install runtime dependencies for cgo and sqlite3
RUN apk add --no-cache libgcc libstdc++

WORKDIR /app

# Copy the binaries
COPY --from=builder /app/myapp .
COPY --from=builder /app/myapp-worker .
# Copy the environment file
COPY --from=builder /app/.env .
# Copy the database
COPY --from=builder /app/test.db .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./myapp"]
