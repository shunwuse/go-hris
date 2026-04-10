# syntax=docker/dockerfile:1.10

# ========================================
# Stage 1: Base build environment
# ========================================
FROM golang:1.26.2-alpine3.23 AS base-builder

# Set the working directory
WORKDIR /src
ENV PATH="/go/bin:${PATH}"

# Install build dependencies
RUN apk add --no-cache ca-certificates

# Install tools
ARG WIRE_VERSION=v0.7.0
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
	--mount=type=cache,id=go-build,target=/root/.cache/go-build \
	go install github.com/google/wire/cmd/wire@${WIRE_VERSION}

# ========================================
# Stage 2: Download dependencies and tools
# ========================================
FROM base-builder AS deps

# Copy dependency files and download
COPY go.mod go.sum ./
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
	--mount=type=cache,id=go-build,target=/root/.cache/go-build \
	go mod download

# ========================================
# Stage 3: Build binaries
# ========================================
FROM deps AS builder

# Copy the rest of the source code
COPY . .

# Accept version as a build argument
ARG VERSION=dev
ARG BUILD_MODE=release

ENV CGO_ENABLED=0
ENV GOOS=linux

# Generate wire files and build the application
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
	--mount=type=cache,id=go-build,target=/root/.cache/go-build \
	wire ./cmd/server && wire ./cmd/worker

RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
	--mount=type=cache,id=go-build,target=/root/.cache/go-build \
	set -eux; \
	if [ "${BUILD_MODE}" = "release" ]; then \
		LDFLAGS="-X github.com/shunwuse/go-hris/internal/infra/app.Version=${VERSION} -s -w -extldflags=-static"; \
		go build -trimpath -ldflags "${LDFLAGS}" -o /out/app ./cmd/server; \
		go build -trimpath -ldflags "${LDFLAGS}" -o /out/app-worker ./cmd/worker; \
	elif [ "${BUILD_MODE}" = "debug" ]; then \
		LDFLAGS="-X github.com/shunwuse/go-hris/internal/infra/app.Version=${VERSION} -extldflags=-static"; \
		go build -gcflags "all=-N -l" -ldflags "${LDFLAGS}" -o /out/app ./cmd/server; \
		go build -gcflags "all=-N -l" -ldflags "${LDFLAGS}" -o /out/app-worker ./cmd/worker; \
	else \
		echo "Unsupported BUILD_MODE: ${BUILD_MODE}"; \
		exit 1; \
	fi

RUN mkdir -p /out/configs /out/logs && \
	cp configs/default.env /out/configs/default.env

# ========================================
# Stage 4: Packaged Artifacts
# ========================================
FROM scratch AS artifacts

COPY --from=builder /out/app /app/app
COPY --from=builder /out/app-worker /app/app-worker
COPY --from=builder /out/configs /app/configs
COPY --from=builder /out/logs /app/logs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# ========================================
# Stage 5: Runtime
# Use this stage for compatibility and debugging tools.
# ========================================
FROM alpine:3.23 AS runtime
# FROM gcr.io/distroless/static-debian12:nonroot AS runtime-distroless
# FROM scratch AS runtime-scratch

WORKDIR /app

COPY --from=artifacts /app/ ./
COPY --from=artifacts /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]
