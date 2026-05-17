#!/bin/sh

repo_root() {
	script_dir=$(CDPATH= cd "$(dirname "$1")" && pwd)
	dirname "$script_dir"
}

default_version() {
	git rev-parse --short HEAD 2>/dev/null || echo dev
}

build_ldflags() {
	version=${VERSION:-$(default_version)}
	printf '%s' "-X github.com/shunwuse/go-hris/internal/infra/app.Version=$version"
}

usage() {
	printf '%s\n' "$1" >&2
	exit 1
}
