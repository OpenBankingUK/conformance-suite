# Functional Conformance Suite - CLI

FCS client line interface provides support to run Open Banking functional conformance test and results is an automated way, typically in a software development pipeline.

## Usage

To check available options with command with out any arguments

```bash
./fcs
```

To use one of the existing templates to generate test cases:

```bash
./fcs run --filename pkg/discovery/templates/ob-v3.1-generic.json --config config.json --output ob-v3.1-generic-testcases.json
```

You can omit `--output` flag and it will write to standard output.
