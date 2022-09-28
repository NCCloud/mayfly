# Contributing Guidelines

Contributions are welcome via GitHub pull requests. This document outlines the process to help get your contribution accepted.

## How to Contribute

1. Fork this repository, develop, and test your changes.
2. Submit a pull request.
3. Make sure all tests are passing.

***NOTE***: In order to make testing and merging of PRs easier, please submit changes for different fixes/features/improvements in a separate PRs.

### Technical Requirements

* Must follow [Golang best practices](https://go.dev/doc/effective_go)
* Must pass CI jobs for linting [golangci-lint](https://github.com/golangci/golangci-lint) tool and unit tests
* All changes require reviews from the responsible organization members before merge.

Once changes have been merged, the release will be done by the responsible organization members.

### Versioning

Versioning should follow [semver](https://semver.org/). Any backwards incompatible changes should be bump the major version and stated in the Release Notes.
