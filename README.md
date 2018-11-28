![Open Banking Logo](https://bitbucket.org/openbankingteam/conformance-suite/raw/99b76db5f60bb4d790d6f32bffae29cbe95a3661/docs/static_files/OBIE_logotype_blue_RGB.PNG)

The **Functional Conformance Suite** is an Open Source test tool provided by [Open Banking](https://www.openbanking.org.uk/). The goal of the suite is to provide an easy and comprehensive tool that enables implementers to test their interfaces/endpoints against the Functional API standard.

## Release Notes 
* * *

### v0.1.0-pre-alpha (30th November 2018)

This **pre-alpha** build introduces some of the concepts to participants before the official release. We are working hard to foster an open and collaborative tool and welcome your feedback to help us develop the best tool.

This release is not intended to be executable so please do **NOT** post bug reports about this version.

Read the full [release notes](docs/releases/v0.1.0-pre-alpha.md)

What's New:

* NEW: [Test Case](docs/test-case-design.md) Design. The conformance suite provides a cross-platform, language agnostic representation of test cases. The primary purpose is to determine whether an implementer resource complies with requirements of specifications and conditions, regulations and standards.
* NEW: [Discovery Configuration](docs/discovery.md) Design: Provides a configurable model that allows an implementer to describe information on endpoint availability, data structure.

To provide feedback, please use the public [issue tracker](https://bitbucket.org/openbankingteam/conformance-suite/issues) or see the [CONTRIBUTING.md](CONTRIBUTING.md).

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

For support on using the suite contact the [Open Banking Service Desk](https://openbanking.atlassian.net/servicedesk/customer/portals). To raise bugs or features please use the [issue tracker](https://bitbucket.org/openbankingteam/conformance-suite/issues).

## Licensing
* * *

This repository is subject to this MIT Open Licence. Please read our [LICENSE.md](LICENSE.md) for more information

## Contributing
* * *
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Useful links
* * *

For more information see the [Open Banking Developer Zone](https://openbanking.atlassian.net/wiki/spaces/DZ/overview).

## Acknowledgments
* * *
