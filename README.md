![Open Banking Logo](https://bitbucket.org/openbankingteam/conformance-suite/raw/99b76db5f60bb4d790d6f32bffae29cbe95a3661/docs/static_files/OBIE_logotype_blue_RGB.PNG)

[![made-with-go](https://img.shields.io/badge/Made%20with-Go-1f425Ff.svg)](https://www.golang.org/)
[![made-with-vue-js](https://img.shields.io/badge/Made%20with-Vue.JS-1f425Ff.svg)](https://vuejs.org/)
[![master](https://img.shields.io/bitbucket/pipelines/openbankingteam/conformance-suite/master.svg)](https://bitbucket.org/openbankingteam/conformance-suite/addon/pipelines/home#!/results/branch/master/page/1)
[![develop](https://img.shields.io/bitbucket/pipelines/openbankingteam/conformance-suite/develop.svg)](https://bitbucket.org/openbankingteam/conformance-suite/addon/pipelines/home#!/results/branch/develop/page/1)
[![Go Reportcard](https://goreportcard.com/badge/bitbucket.org/openbankingteam/conformance-suite)](https://goreportcard.com/report/bitbucket.org/openbankingteam/conformance-suite)

---

The **Functional Conformance Tool** is an Open Source test tool provided by [Open Banking](https://www.openbanking.org.uk/). The goal of the suite is to provide an easy and comprehensive tool that enables implementers to test interfaces and data endpoints against the Functional API standard.

The supporting documentation assumes technical understanding of the Open Banking ecosystem. An introduction to the concepts is available via the [Open Banking Website](https://www.openbanking.org.uk/).

To provide feedback, please use the public [issue tracker](https://bitbucket.org/openbankingteam/conformance-suite/issues) or see the [CONTRIBUTING.md](CONTRIBUTING.md).

## Release Notes 
* * *

# Release v1.6.7 (24th February 2021)

The release is called **v1.6.7**, an update to support [additional roadmap tests](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1564083551/OBIE+Functional+Conformance+Tool+Roadmap) and contains minor fixes.

[Full Release Notes](https://bitbucket.org/openbankingteam/conformance-suite/src/v1.6.7/docs/releases/v1.6.7.md)
(v1.6.7.md)

## Quickstart
* * *

Pull and run the latest (stable) tagged Docker image:

    > docker run --rm -it -p 8443:8443 "openbanking/conformance-suite:v1.6.7"

[See Setup Guide](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/docs/setup-guide.md)

### Prerequisites

The tool is compatible with the Open Banking UK R/W specification versions: 3.1.0, 3.1.1, 3.1.2, 3.1.3, 3.1.4, 3.1.5, 3.1.6.

In order to run a container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

## Support
* * *

For support on using the suite use the [Open Banking Help Centre](https://openbanking.atlassian.net/servicedesk/customer/portals). To raise bugs or features please use the [issue tracker](https://bitbucket.org/openbankingteam/conformance-suite/issues).

## Licensing
* * *

This repository is subject to this MIT Open Licence. Please read our [LICENSE.md](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/LICENSE.md) for more information

## Contributing
* * *
Please read [CONTRIBUTING.md](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## Useful links
* * *

* [Docker Conformance Tool](https://hub.docker.com/r/openbanking/conformance-suite/)
* [Open Banking Developer Zone](https://openbanking.atlassian.net/wiki/spaces/DZ/overview)
* [All Release Notes](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/docs/releases/releases.md)
