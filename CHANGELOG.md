# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-27

### Added

- JSON Schema parsing (JSON and YAML formats)
- OpenAPI 3.x specification parsing and schema extraction
- Schema diff engine with 17 breaking change detection rules
- `$ref` pointer resolution with circular reference handling
- CLI with `check` and `rules` commands
- Text output with optional ANSI color
- JSON output (compact and pretty-printed)
- SARIF 2.1.0 output for CI integration
- Configurable fail conditions (`--fail-on breaking|warning|any|none`)
- Auto-detection of input format (JSON Schema vs OpenAPI)
- Exit code support for CI gating (0 = compatible, 1 = breaking, 2 = error)
- 108 tests covering all packages
