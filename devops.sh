#!/usr/bin/env bash

export CONTROLLER_GEN_VERSION="v0.15.0"
export GOLANGCI_LINT_VERSION="v1.59.1"
export MOCKERY_GEN_VERSION="v2.44.1"
export GOFUMPT_VERSION="v0.6.0"
export TESTENV_VERSION="1.25.x!"

prerequisites() {
  if [[ "$(controller-gen --version 2>&1)" != *"$CONTROLLER_GEN_VERSION"* ]]; then
    go install sigs.k8s.io/controller-tools/cmd/controller-gen@"${CONTROLLER_GEN_VERSION}"
  fi
  if [[ "$(golangci-lint --version 2>&1)" != *"$GOLANGCI_LINT_VERSION"* ]]; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@"${GOLANGCI_LINT_VERSION}"
  fi
  if [[ "$(mockery --version 2>&1)" != *"$MOCKERY_GEN_VERSION"* ]]; then
    go install github.com/vektra/mockery/v2@"${MOCKERY_GEN_VERSION}"
  fi
  if [[ "$(gofumpt --version 2>&1)" != *"$GOFUMPT_VERSION"* ]]; then
     go install mvdan.cc/gofumpt@"${GOFUMPT_VERSION}"
  fi
  if ! command -v crd-ref-docs &>/dev/null; then
    go install github.com/elastic/crd-ref-docs@latest
  fi
  if ! command -v setup-envtest &>/dev/null; then
    go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
  fi
}

lint() {
  gofumpt -l -w .
  golangci-lint run --timeout=10m
}

generate() {
  rm -rf deploy/crds
  controller-gen object paths="./..."
  controller-gen crd paths="./..." output:dir=deploy/crds
  crd-ref-docs --source-path=./pkg/apis --config .apidoc.yaml --renderer markdown --output-path=./docs/api.md
  mockery
}

install() {
  kubectl apply -f deploy/crds
}

prepare_envtest() {
  mkdir -p .envtest/crds
  mkdir -p .envtest/bins
  cp -rf "$(setup-envtest use $TESTENV_VERSION -p path)"/* .envtest/bins/
}

test() {
  go test -v -coverpkg=./... ./...
}

prerequisites

"$@"