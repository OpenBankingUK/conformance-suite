![Open Banking Logo](https://raw.githubusercontent.com/OpenBankingUK/conformance-suite/develop/docs/static_files/OBIE_logotype_blue_RGB.PNG)

[![made-with-go](https://img.shields.io/badge/Made%20with-Go-1f425Ff.svg)](https://www.golang.org/)
[![made-with-vue-js](https://img.shields.io/badge/Made%20with-Vue.JS-1f425Ff.svg)](https://vuejs.org/)
[![master](https://img.shields.io/github/checks-status/OpenBankingUK/conformance-suite/master.svg)](https://github.com/OpenBankingUK/conformance-suite/actions?query=branch%3Amaster+)
[![develop](https://img.shields.io/github/checks-status/OpenBankingUK/conformance-suite/develop.svg)](https://github.com/OpenBankingUK/conformance-suite/actions?query=branch%3Adevelop)
[![Go Reportcard](https://goreportcard.com/badge/github.com/OpenBankingUK/conformance-suite)](https://goreportcard.com/report/github.com/OpenBankingUK/conformance-suite)

---

The **Functional Conformance Tool** is an Open Source test tool provided by [Open Banking](https://www.openbanking.org.uk/). The goal of the suite is to provide an easy and comprehensive tool that enables implementers to test interfaces and data endpoints against the Functional API standard.

The supporting documentation assumes a technical understanding of the Open Banking ecosystem. An introduction to the concepts is available via the [Open Banking Website](https://www.openbanking.org.uk/).

To provide feedback, please see the [CONTRIBUTING.md](CONTRIBUTING.md).

## Release Notes
* * *

# Release v1.9.1 (23rd October 2024)

The release is called **v1.9.1**, an update to add minor bugfixes and security updates, such as xss security headers, implements the FCS re-run functionality, VRPType 3.1.10/3.1.11 support, nbf jwt token field and many more.

[Full Release Notes](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/releases/v1.9.1.md)


---
**Download**:
`docker run --rm -it -p 127.0.0.1:8443:8443 "openbanking/conformance-suite:v1.9.1"` |
[DockerHub](https://hub.docker.com/r/openbanking/conformance-suite) |
[Setup Guide](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/setup-guide.md)
---


## Version table

| Release       | Standard version  |
| ------------- |:-----------------:|
| v1.9.0        | v4.0.0            |
| v1.7.6        | v3.1.11           |
| v1.7.0        | v3.1.10           |
| v1.6.12       | v3.1.9            |


## Quickstart
* * *

Pull and run the latest (stable) tagged Docker image:

    > docker run --rm -it -p 127.0.0.1:8443:8443 "openbanking/conformance-suite:v1.9.1"

or

    > docker run --rm -it -p 8443:8443 "openbanking/conformance-suite:v1.9.1"

[See Setup Guide](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/setup-guide.md)

### Prerequisites

The tool is compatible with the Open Banking UK R/W specification versions: 3.1.0, 3.1.1, 3.1.2, 3.1.3, 3.1.4, 3.1.5, 3.1.6, 3.1.7, 3.1.8, 3.1.9, 3.1.10, 3.1.11, 4.0.0.

In order to run a container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

## Advanced Logging

There is the ability to enhance logging by setting any of the following env variables

- LOG_HTTP_TRACE=true:  Enables detailed HTTP trace logging. When set to true, the application will log all HTTP requests and responses, including headers and body content. This can be useful for debugging and monitoring HTTP interactions, but may expose sensitive information in the logs. Use with caution in production environments.
- LOG_LEVEL=debug: sets the logging level to debug, providing detailed information about the application's operation. This is useful for diagnosing issues and understanding the application's behavior in detail.
- LOG_TRACER=true: enables detailed tracing of the application's execution. This setting provides in-depth information about the application's internal processes, which can be useful for debugging complex issues. Use with caution as it may generate a large volume of log data.
- EXPORT_LOG_FILE=true: sets the log output to a file and attaches the log file to the run export zip archive

## Support
* * *

For support on using the suite use the [Open Banking Help Centre](https://openbanking.atlassian.net/servicedesk/customer/portals).

## Licensing
* * *

This repository is subject to this MIT Open Licence. Please read our [LICENSE.md](https://github.com/OpenBankingUK/conformance-suite/blob/develop/LICENSE.md) for more information

## Contributing
* * *
Please read [CONTRIBUTING.md](https://github.com/OpenBankingUK/conformance-suite/blob/develop/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## Useful links
* * *

* [Docker Conformance Tool](https://hub.docker.com/r/openbanking/conformance-suite/)
* [Open Banking Developer Zone](https://openbanking.atlassian.net/wiki/spaces/DZ/overview)
* [All Release Notes](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/releases/releases.md)
