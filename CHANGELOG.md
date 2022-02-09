# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.6.0] - 2022-02-09

### Added

- Add `priorityClassName: "system-node-critical"` to Daemonset to give higher priority during scheduling.

## [0.5.0] - 2021-12-06

### Added

- Added basic prometheus exporter.

## [0.4.0] - 2021-04-13

### Added

- React to `Preempt`, `Reboot` and `Redeploy` events to drain nodes properly.

### Change

- Use the `NotBefore` field from the event to calculate drain timeout.

## [0.3.0] - 2021-03-19

### Fixed

- Ensure to wait long enough when draining a node before considering the node drained.

### Changed

- Change drain timeout to 15 minutes.

## [0.2.2] - 2021-02-23

### Fixed

- Disable helm hook for new installations of the chart.

## [0.2.1] - 2021-02-22

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

[Unreleased]: https://github.com/giantswarm/giantswarm/compare/v0.6.0...HEAD
[0.6.0]: https://github.com/giantswarm/giantswarm/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.2.2...v0.3.0
[0.2.2]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/azure-scheduled-events/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/giantswarm/azure-operator/releases/tag/v0.0.1
