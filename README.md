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

# Release v1.8.1 (12th August 2024)

The release is called **v1.8.1**, an update to add preliminary support to Open Banking Read/Write APIs v4.0.0 specification, including some UX/UI improvements for new default manifests compatible with v4.0.0, new swagger and minor bugfixes.

[Full Release Notes](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/releases/v1.8.1.md)


---
**Download**:
`docker run --rm -it -p 127.0.0.1:8443:8443 "openbanking/conformance-suite:v1.8.1"` |
[DockerHub](https://hub.docker.com/r/openbanking/conformance-suite) |
[Setup Guide](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/setup-guide.md)
---


## Version table

| Release       | Standard version  |
| ------------- |:-----------------:|
| v1.8.1        | v4.0.0            |
| v1.7.2        | v3.1.11           |
| v1.7.0        | v3.1.10           |
| v1.6.12       | v3.1.9            |


## Quickstart
* * *

Pull and run the latest (stable) tagged Docker image:

    > docker run --rm -it -p 127.0.0.1:8443:8443 "openbanking/conformance-suite:v1.8.0"

or

    > docker run --rm -it -p 8443:8443 "openbanking/conformance-suite:v1.8.0"

[See Setup Guide](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/setup-guide.md)

### Prerequisites

The tool is compatible with the Open Banking UK R/W specification versions: 3.1.0, 3.1.1, 3.1.2, 3.1.3, 3.1.4, 3.1.5, 3.1.6, 3.1.7, 3.1.8, 3.1.9, 3.1.10, 3.1.11, 4.0.0.

In order to run a container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

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
