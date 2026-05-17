#!/bin/sh
set -eu

. "$(dirname "$0")/lib.sh"

cd "$(repo_root "$0")"

go run github.com/google/wire/cmd/wire@v0.7.0 ./cmd/server
go run github.com/google/wire/cmd/wire@v0.7.0 ./cmd/worker
go mod tidy
