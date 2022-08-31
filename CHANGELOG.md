# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
 - Build binaries for darwin and windows (not tested, please do report if you encounter any issues!)
### Changed
 - Upgraded to Go 1.19
 - Ignore OpenMetrics `_created` metrics (for more info see [this](https://github.com/prometheus/prometheus/issues/6541) issue)
 - Handle metrics ending with `_time|_time_seconds|_timestamp|_timestamp_seconds` as timestamps in seconds and subtract them from `time()`

## [1.1.0] - 2022-08-12
### Added
 - Short flags for most common flags, see `--help` or [readme](https://github.com/FUSAKLA/autograf#how-to-use)
### Fixed
 - Correctly handle metric selectors if no `--selector` is set
### Changed
 - Metric with unknown type is now visualized as gauge in time series panel as "best effort"

## [1.0.1] - 2022-07-12
 - Fixed default folder name

## [1.0.0] - 2022-07-12
Initial release
