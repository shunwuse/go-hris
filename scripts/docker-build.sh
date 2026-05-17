#!/bin/sh
set -eu

. "$(dirname "$0")/lib.sh"

cd "$(repo_root "$0")"

PLATFORM=${PLATFORM:-linux/amd64}
DOCKER_IMAGE=${DOCKER_IMAGE:-go-hris:latest}
DOCKERFILE=${DOCKERFILE:-build/package/Dockerfile}
BUILD_MODE=${BUILD_MODE:-release} # debug or release

if [ -z "${VERSION:-}" ]; then
	VERSION=$(default_version)
fi

exec docker buildx build \
	--platform "$PLATFORM" \
	--build-arg BUILD_MODE="$BUILD_MODE" \
	--build-arg VERSION="$VERSION" \
	-t "$DOCKER_IMAGE" \
	--load \
	-f "$DOCKERFILE" \
	.
