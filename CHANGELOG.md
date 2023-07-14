# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.2.0] - 2023-07-14
### Added
 - New config option `prometheus_bearer_token` to allow authentication with Prometheus using a bearer token
### Changed
 - Upgraded to go 1.20
 - Upgraded all dependencies to newest versions

## [2.1.0] - 2023-05-06
### Added
 - New flag `--ignore-config | -i` to ignore any config file found
 - Updated demo example to newest grafana 9.5.1 and verified all works
 - Updated all dependencies to newest versions
### Fixed
 - Fix heatmap panel to not calculate the histogram buckets

## [2.0.1] - 2022-11-08
### Fixed
 - Fix heatmap panel query format to be valid

## [2.0.0] - 2022-11-06
### Added
 - Build binaries for darwin and windows (not tested, please do report if you encounter any issues!)
 - Config is now loaded also from location `~/.config.autograf.json` (Which config is used is printed to stderr to ease debugging)
 - Support more metric units (amper, hertz, volt, ...)
 - Experimental support for customization of the generated panels using the metric HELP text, see [the docs](./README.md#panel-config-customization-experimental)
 - Added demo in the [`./examples/demo`](./examples/demo/)
 - `--version` flag to print out Autograf version
### Changed
 - Optimized for the newest Grafana releases (9.2+)
 - Use the new heatmap Grafana panel
 - Switch range selector from `3m` to special Grafana variable `$__rate_interval`
 - Upgraded to Go 1.19
 - Ignore OpenMetrics `_created` metrics (for more info see [this](https://github.com/prometheus/prometheus/issues/6541) issue)
 - Handle metrics ending with `_time|_time_seconds|_timestamp|_timestamp_seconds` as timestamps in seconds and subtract them from `time()`
 - Improve layout of generated panels
 - Add metric names in the name of the row for better visibility

## [1.1.0] - 2022-08-12
### Added
 - Short flags for most common flags, see `--help` or [readme](https://github.com/fusakla/autograf#how-to-use)
### Fixed
 - Correctly handle metric selectors if no `--selector` is set
### Changed
 - Metric with unknown type is now visualized as gauge in time series panel as "best effort"

## [1.0.1] - 2022-07-12
 - Fixed default folder name

## [1.0.0] - 2022-07-12
Initial release
