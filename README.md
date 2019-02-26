![Open Banking Logo](https://bitbucket.org/openbankingteam/conformance-suite/raw/99b76db5f60bb4d790d6f32bffae29cbe95a3661/docs/static_files/OBIE_logotype_blue_RGB.PNG)

[![made-with-go](https://img.shields.io/badge/Made%20with-Go-1f425Ff.svg)](https://www.golang.org/)
[![made-with-vue-js](https://img.shields.io/badge/Made%20with-Vue.JS-1f425Ff.svg)](https://vuejs.org/)
[![master](https://img.shields.io/bitbucket/pipelines/openbankingteam/conformance-suite/master.svg)](https://bitbucket.org/openbankingteam/conformance-suite/addon/pipelines/home#!/results/branch/master/page/1)
[![develop](https://img.shields.io/bitbucket/pipelines/openbankingteam/conformance-suite/develop.svg)](https://bitbucket.org/openbankingteam/conformance-suite/addon/pipelines/home#!/results/branch/develop/page/1)
[![Go Reportcard](https://goreportcard.com/badge/bitbucket.org/openbankingteam/conformance-suite)](https://goreportcard.com/report/bitbucket.org/openbankingteam/conformance-suite)

The **Functional Conformance Suite** is an Open Source test tool provided by [Open Banking](https://www.openbanking.org.uk/). The goal of the suite is to provide an easy and comprehensive tool that enables implementers to test interfaces and data endpoints against the Functional API standard.

The supporting documentation assumes technical understanding of the Open Banking ecosystem. An introduction to the concepts is available via the [Open Banking Website](https://www.openbanking.org.uk/).

To provide feedback, please use the public [issue tracker](https://bitbucket.org/openbankingteam/conformance-suite/issues) or see the [CONTRIBUTING.md](CONTRIBUTING.md).

## Release Notes 
* * *

### v1.0.0-beta (19th February 2019)

This **v1.0.0-beta** release introduces manual PSU consent (hybrid flow) and a example Discovery Template to demonstrate the hybrid flow for Ozone Model Bank for v3.1 of the OBIE Accounts and Transactions specifications.

N.B. This release is not intended to be run in production.

Read the full [release notes](docs/releases/v1.0.0-beta.md) (v1.0.0-beta.md)

### v0.2.0-alpha (11th January 2019)

This **v0.2.0-alpha** release introduces an example Discovery Template to demonstrate the headless flow for Ozone Model Bank for v3.0 of the OBIE  Accounts and Transactions specifications.

N.B. This release is not intended to be run in sandbox or production so bug reports about this version should **NOT** be posted.

Read the full [release notes](docs/releases/v0.2.0-alpha.md) (v0.2.0-alpha.md)

### v0.1.0-pre-alpha (30th November 2018)

This **pre-alpha** build introduces some of the concepts to implementers before the official release. The aim is to foster an open and collaborative tool and support feedback to help develop the best possible tool.

N.B. This release is not intended to be executable so bug reports about this version should **NOT** be posted.

Read the full [release notes](docs/releases/v0.1.0-pre-alpha.md) (v0.1.0-pre-alpha.md)

What's New:

* NEW: [Test Case Design](docs/test-case-design.md) : The functional conformance suite provides a cross-platform, language agnostic representation of test cases. The primary purpose is to determine whether an implementer resource complies with requirements of specifications and conditions, regulations and standards.
* NEW: [Discovery Design](docs/discovery.md): Provides a configurable model that allows an implementer to describe information on endpoint availability and data structure.

## Quickstart
* * *

Clone:

    
    > git clone git@bitbucket.org:openbankingteam/conformance-suite.git && cd conformance-suite


Run the Makefile:


    > make run_image


You can also pull & run the latest Docker image:


    > docker run openbanking/conformance-suite


### Prerequisites

In order to run this container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

## Support
* * *

For support on using the suite use the [Open Banking Help Centre](https://openbanking.atlassian.net/servicedesk/customer/portals). To raise bugs or features please use the [issue tracker](https://bitbucket.org/openbankingteam/conformance-suite/issues).

## Licensing
* * *

This repository is subject to this MIT Open Licence. Please read our [LICENSE.md](LICENSE.md) for more information

## Contributing
* * *
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## Useful links
* * *

* Docker Hub [Conformance Suite](https://hub.docker.com/r/openbanking/conformance-suite/)
* For more information see the [Open Banking Developer Zone](https://openbanking.atlassian.net/wiki/spaces/DZ/overview).
