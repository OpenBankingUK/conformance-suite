# Functional Conformance Suite - CLI

FCS client line interface provides support to run Open Banking functional conformance test and results is an automated way, typically in a software development pipeline.

## Configuration

Preferred configuration type is using environment variables following [12 factor recommendations](https://12factor.net/config).

Configuration variable are prefixed with `fcs`.

Comprehensive list of all configuration available:

```bash
FCS_WELCOME=Placeholder
```

## Usage

To check available options with command with out any arguments

```bash
./fcs
```

To use one of the existing templates to generate test cases:

```bash
./fcs generate --filename pkg/discovery/templates/ob-v3.1-generic.json --output ob-v3.1-generic-testcases.json
```

You can omit `--output` flag and it will write to standard output.
