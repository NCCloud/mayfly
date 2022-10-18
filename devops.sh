#!/usr/bin/env bash

export BLUE="\033[0;34m\033[1m"
export NC="\033[0m"

export GOLANGCI_LINT_VERSION="v1.47.2"

prerequisites() {
  if ! command -v golangci-lint &>/dev/null; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@"${GOLANGCI_LINT_VERSION}"
  fi
  if ! command -v gofumpt &>/dev/null; then
    go install mvdan.cc/gofumpt@latest
  fi
}

lint() {
  gofumpt -l -w .
  golangci-lint run --timeout=10m
}

test() {
  go test -v ./...
}

prerequisites

"$@"