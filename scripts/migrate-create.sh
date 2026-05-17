#!/bin/sh
set -eu

. "$(dirname "$0")/lib.sh"

cd "$(repo_root "$0")"

name=${1:-${name:-}}
if [ -z "$name" ]; then
	usage "Usage: $(basename "$0") <migration_name>"
fi

exec go run github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0 create -ext sql -dir ./migrations -seq "$name"
