#!/bin/sh
set -eu

. "$(dirname "$0")/lib.sh"

cd "$(repo_root "$0")"

if [ -z "${LDFLAGS:-}" ]; then
	LDFLAGS=$(build_ldflags)
fi

exec go run -ldflags "$LDFLAGS" ./cmd/worker
