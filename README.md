<div align="center">

# schemadiff

### Detect breaking schema changes before they ship

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Tests](https://img.shields.io/badge/Tests-108+-success?style=for-the-badge)](https://github.com/JSLEEKR/schemadiff)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![SARIF](https://img.shields.io/badge/SARIF-2.1.0-orange?style=for-the-badge)](https://sarifweb.azurewebsites.net/)

<br/>

**Catch backward-incompatible changes in JSON Schema and OpenAPI before they reach production.**

schemadiff compares two schema versions and reports every breaking change, warning,
and informational change. Built for CI pipelines with SARIF output, configurable
exit codes, and zero dependencies beyond the Go standard library + YAML/Cobra.

</div>

---

## Why This Exists

Schema changes are invisible breakage. A renamed field, a tightened constraint, a removed
enum value --- these slip through code review because the diff looks harmless. Consumers
discover the break at runtime, in production, on a Friday night.

**schemadiff** makes schema compatibility a CI gate. Run it on every PR that touches a
schema file. If the change is backward-incompatible, the build fails before merge.

**The problem:**
- API schemas evolve across teams and services
- JSON Schema and OpenAPI specs change without backward-compatibility checks
- Breaking changes surface at runtime, not at review time
- Existing tools are either too slow, too complex, or don't support CI output formats

**The solution:**
- Single binary, zero config, instant feedback
- 17 built-in rules covering all common breaking change patterns
- SARIF 2.1.0 output for GitHub Code Scanning integration
- Exit codes for CI/CD gating

---

## Quick Start

### Install

```bash
go install github.com/JSLEEKR/schemadiff@latest
```

### Basic Usage

```bash
# Compare two JSON Schema files
schemadiff check old.json new.json

# Compare OpenAPI specs
schemadiff check old-api.yaml new-api.yaml --format openapi

# JSON output
schemadiff check old.json new.json --output json

# SARIF output for CI
schemadiff check old.json new.json --output sarif

# Fail on any change (not just breaking)
schemadiff check old.json new.json --fail-on any

# Never fail (report only)
schemadiff check old.json new.json --fail-on none

# List all rules
schemadiff rules
```

---

## Breaking Change Rules

schemadiff detects 17 types of schema changes across three severity levels:

### BREAKING (exit code 1)

| Rule | Description |
|------|-------------|
| `required-field-added` | Adding a required field breaks existing consumers |
| `property-removed` | Removing a property breaks consumers that depend on it |
| `type-changed` | Changing a field's type breaks type expectations |
| `enum-value-removed` | Removing an enum value breaks consumers using it |
| `min-length-increased` | Tighter minLength rejects previously valid strings |
| `max-length-decreased` | Tighter maxLength rejects previously valid strings |
| `minimum-increased` | Higher minimum rejects previously valid numbers |
| `maximum-decreased` | Lower maximum rejects previously valid numbers |
| `min-items-increased` | Higher minItems rejects previously valid arrays |
| `max-items-decreased` | Lower maxItems rejects previously valid arrays |
| `pattern-changed` | Changed pattern rejects previously valid strings |
| `additional-properties-restricted` | Restricting extra properties rejects valid payloads |
| `items-schema-changed` | Changed array items schema breaks consumers |

### WARNING

| Rule | Description |
|------|-------------|
| `format-changed` | Format change may break consumers with format validation |

### INFO (non-breaking)

| Rule | Description |
|------|-------------|
| `required-field-removed` | Making a field optional is backward-compatible |
| `property-added` | Adding an optional property is backward-compatible |
| `enum-value-added` | Adding an enum value is backward-compatible |

---

## Output Formats

### Text (default)

```
## Schema
==========

  [X] BREAKING  properties.name.type
     Type changed from "string" to "number"
     Old: string -> New: number

  [i] INFO      properties.phone
     Property "phone" added

  Summary: 1 breaking, 0 warning, 1 info (2 total)
```

### JSON

```json
{
  "version": "1.0.0",
  "result": {
    "changes": [
      {
        "path": "properties.name.type",
        "type": "type-changed",
        "severity": 2,
        "message": "Type changed from \"string\" to \"number\"",
        "old_value": "string",
        "new_value": "number"
      }
    ],
    "has_breaking": true,
    "summary": {
      "breaking": 1,
      "warning": 0,
      "info": 0,
      "total": 1
    }
  },
  "summary": {
    "total_schemas": 1,
    "total_changes": 1,
    "breaking": 1,
    "warning": 0,
    "info": 0,
    "has_breaking": true
  }
}
```

### SARIF 2.1.0

SARIF output integrates with GitHub Code Scanning, Azure DevOps, and other SARIF-compatible tools:

```json
{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "schemadiff",
          "version": "1.0.0",
          "rules": [...]
        }
      },
      "results": [...]
    }
  ]
}
```

---

## CI Integration

### GitHub Actions

```yaml
name: Schema Check
on:
  pull_request:
    paths:
      - 'schemas/**'
      - 'api/**'

jobs:
  schema-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Install schemadiff
        run: go install github.com/JSLEEKR/schemadiff@latest

      - name: Check for breaking changes
        run: |
          git show HEAD~1:schemas/api.json > /tmp/old.json
          schemadiff check /tmp/old.json schemas/api.json --output sarif > results.sarif

      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif
```

### GitLab CI

```yaml
schema-check:
  stage: test
  image: golang:1.26
  script:
    - go install github.com/JSLEEKR/schemadiff@latest
    - git show HEAD~1:schemas/api.json > /tmp/old.json
    - schemadiff check /tmp/old.json schemas/api.json
  only:
    changes:
      - schemas/**
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

SCHEMA_FILES=$(git diff --cached --name-only --diff-filter=M -- '*.json' | grep -E 'schema|api')
for file in $SCHEMA_FILES; do
  OLD=$(git show HEAD:"$file" 2>/dev/null)
  if [ -n "$OLD" ]; then
    echo "$OLD" > /tmp/old_schema.json
    schemadiff check /tmp/old_schema.json "$file" --fail-on breaking
    if [ $? -ne 0 ]; then
      echo "Breaking schema change detected in $file"
      exit 1
    fi
  fi
done
```

---

## OpenAPI Support

schemadiff extracts and compares schemas from all locations in an OpenAPI 3.x spec:

- **Component schemas** (`components.schemas.*`)
- **Request bodies** (`paths.*.*.requestBody.content.*.schema`)
- **Response bodies** (`paths.*.*.responses.*.content.*.schema`)
- **Parameters** (`paths.*.*.parameters[*].schema`)

```bash
# Auto-detected from spec content
schemadiff check old-api.yaml new-api.yaml

# Explicit format
schemadiff check old-api.json new-api.json --format openapi
```

The output shows changes per schema location:

```
## POST /users request (application/json)
=============================================

  [X] BREAKING  required
     Field "phone" added to required

## GET /users 200 response (application/json)
===============================================

  [X] BREAKING  items.properties.id.type
     Type changed from "string" to "integer"

Overall: 2 breaking, 0 warning, 0 info (2 total changes across 2 schemas)
```

---

## CLI Reference

### `schemadiff check <old> <new>`

Compare two schema files and report changes.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `text` | Output format: `text`, `json`, `sarif` |
| `--format` | `-f` | `auto` | Input format: `auto`, `jsonschema`, `openapi` |
| `--fail-on` | | `breaking` | Fail condition: `breaking`, `warning`, `any`, `none` |
| `--no-color` | | `false` | Disable ANSI colored output |

**Exit codes:**
- `0` --- No breaking changes (or fail condition not met)
- `1` --- Breaking changes detected (or fail condition met)
- `2` --- Invalid input or error

### `schemadiff rules`

List all 17 breaking change detection rules with their severity and description.

### `schemadiff --version`

Print the version number.

---

## Architecture

```
schemadiff/
├── main.go                    # Entry point
├── cmd/
│   ├── root.go                # Cobra root command
│   ├── check.go               # check subcommand
│   └── rules.go               # rules subcommand
├── pkg/
│   ├── schema/
│   │   ├── types.go           # Core types (Schema, Change, DiffResult)
│   │   ├── parser.go          # JSON/YAML parser with $ref resolution
│   │   ├── differ.go          # Schema diff engine
│   │   └── rules.go           # Breaking change rule definitions
│   ├── openapi/
│   │   ├── parser.go          # OpenAPI 3.x parser
│   │   └── extractor.go       # Schema extraction from OpenAPI specs
│   └── report/
│       ├── text.go            # Human-readable text output
│       ├── json.go            # JSON output
│       └── sarif.go           # SARIF 2.1.0 output
└── internal/
    ├── integration_test.go    # End-to-end tests
    └── testdata/              # Test fixtures
```

### Design Decisions

1. **Recursive diff engine**: Handles nested objects, arrays, and `$ref` pointers
   with circular reference detection.

2. **Severity-based classification**: Every change is classified as BREAKING,
   WARNING, or INFO based on backward-compatibility impact.

3. **Format-agnostic core**: The diff engine works on `Schema` structs.
   Parsers (JSON Schema, OpenAPI) normalize input to this common type.

4. **CI-first design**: SARIF output, exit codes, and `--fail-on` flags
   are first-class features, not afterthoughts.

---

## Supported Schema Features

| Feature | Support |
|---------|---------|
| JSON Schema Draft 4/6/7 | Properties, required, type, enum, constraints |
| `$ref` resolution | Local definitions with circular ref detection |
| `allOf`/`anyOf`/`oneOf` | Parsed but not yet diffed (tracked for v2) |
| OpenAPI 3.0/3.1 | Components, request/response bodies, parameters |
| YAML input | Full support for `.yaml`/`.yml` files |
| Nested objects | Recursive property comparison |
| Array items | Items schema diffing |
| String constraints | minLength, maxLength, pattern, format |
| Numeric constraints | minimum, maximum |
| Array constraints | minItems, maxItems |
| additionalProperties | Restriction detection |

---

## Performance

schemadiff is designed for speed in CI pipelines:

- **Single binary**: No runtime dependencies, no Docker required
- **Instant startup**: Go binary, no JVM or interpreter warmup
- **Linear complexity**: O(n) where n is the number of schema properties
- **Memory efficient**: Streaming output, no full AST in memory

Typical execution time: **< 50ms** for schemas with hundreds of properties.

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new rules or features
4. Run `go test ./...` to verify
5. Submit a pull request

### Adding a New Rule

1. Add the `ChangeType` constant to `pkg/schema/types.go`
2. Add the detection logic to `pkg/schema/differ.go`
3. Add the description to `pkg/schema/rules.go`
4. Add tests to `pkg/schema/differ_test.go`
5. Add the SARIF rule mapping to `pkg/report/sarif.go`

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

<div align="center">

**schemadiff** --- Because breaking changes should break the build, not production.

Built with Go. Designed for CI. Ships as a single binary.

</div>
