# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.3] - 2024-01-12

### Bug fix

* Fixes `LICENSE` to be recognizable by pkg.go.dev. See [#12](https://github.com/ccin2p3/go-freeipa/pull/12)

## [1.1.2] - 2023-08-21

### Bug fix

* Fixes _filescriptor leak_ on `login` by @gladkovandrey. See [#11](https://github.com/ccin2p3/go-freeipa/pull/11)

## [1.1.1] - 2022-10-06

### Added

* Add a complete example on how to list hosts in the freeipa domain

### Bug fix

* Fixes [#7](https://github.com/ccin2p3/go-freeipa/issues/7) by making the host `ManagedbyHost` parameter optional

## [1.1.0] - 2022-04-05

### Changed

- `LICENSE` change for new code ownership
- Take over the maintenance of this library. [Upstream project is unmaintained](https://github.com/tehwalris/go-freeipa#unmaintained)
- Rename the package name from `github.com/tehwalris/go-freeipa` to `github.com/ccin2p3/go-freeipa`
