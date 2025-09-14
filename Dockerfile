# Description: Dockerfile for building the application
FROM golang:1.25.1 AS builder

# copy the source code
COPY . /app

# set the working directory
WORKDIR /app

# install the migrate tool
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# run the migrations
RUN make migrate-up

# download the dependencies
RUN go mod download

# install wire
RUN go install github.com/google/wire/cmd/wire@latest
# generate the wire files
RUN make wire

# build the application
RUN GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o myapp ./cmd/server

# create a new image
FROM alpine:3.20

# Install runtime dependencies for cgo and sqlite3
RUN apk add --no-cache libgcc libstdc++

# copy the binary
COPY --from=builder /app/myapp .
# copy the environment file
COPY --from=builder /app/.env .
# copy the database
COPY --from=builder /app/test.db .

# expose the port
EXPOSE 8080

# run the application
CMD ["./myapp"]
