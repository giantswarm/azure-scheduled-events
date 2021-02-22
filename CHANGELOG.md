# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Fixed labels in the Helm chart.

## [0.2.0] - 2021-02-22

### Added

- Remove the `Node` from Kubernetes API server right before approving the termination event.

### Fixed

- Keep looping on events if one loop errors out.

## [0.1.1] - 2021-01-27

### Changed

- Renamed helm chart with `-app` suffix.

## [0.1.0] - 2021-01-25

## [0.0.1] - 2021-01-20

### Added

- Initial release.

[Unreleased]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/giantswarm/azure-operator/releases/tag/v0.0.1
