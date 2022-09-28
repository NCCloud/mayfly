#!/usr/bin/env bash

export GOLANGCI_LINT_VERSION="v1.50.1"

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