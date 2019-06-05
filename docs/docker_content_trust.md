## Docker Content Trust - Signing Docker Images

Docker Content Trust *(DCT)* allows Open Banking to sign docker images before they get published to the Docker Hub registry. These signatures allow client-side or runtime verification of the integrity and publisher of specific image tags

## How Open Banking Sign Images With Docker Content Trust

Open Banking sign images with the `$ docker trust` command. Which is built on top of the [Notary cli](https://docs.docker.com/notary/getting_started/). Notary is a tool for publishing and managing trusted collections of content. Publishers can digitally sign collections and consumers can verify integrity and origin of content.

### Generate Docker Signing Keys

To sign a Docker image we delegation a key pair. These keys can be generated locally using `$ docker trust key generate`

```sh
$ docker trust key generate jeff
Generating key for jeff...
Enter passphrase for new jeff key with ID 9deed25: 
Repeat passphrase for new jeff key with ID 9deed25: 
Successfully generated and loaded private key. Corresponding public key available: /home/ubuntu/workspace/conformance-suite/jeff.pub
```
Next we will need to add the delegation public key to the Notary server; this is specific to a particular image repository in Notary.

```sh
$ docker trust signer add --key jeff.pub jeff docker.io/openbanking/conformance-suite
Adding signer "jeff" to docker.io/openbanking/conformance-suite...
Enter passphrase for new repository key with ID 10b5e94: 
```

### Enable Docker Image Signing And Publish Image

Once the keys have been imported an image can be pushed with the $ docker push command, by exporting the DCT environmental variable.

```sh
$ export DOCKER_CONTENT_TRUST=1

$ docker push docker.io/openbanking/conformance-suite:v1.0.0
The push refers to repository [docker.io/openbanking/conformance-suite:v1.0.0]
7bff100f35cb: Pushed 
v1.0.0: digest: sha256:3d2e482b82608d153a374df3357c0291589a61cc194ec4a9ca2381073a17f58e size: 528
Signing and pushing trust metadata
Enter passphrase for signer key with ID 8ae710e: 
Successfully signed docker.io/openbanking/conformance-suite:v1.0.0
```

### How to Verify Image

Remote trust data for a tag or a repository can be viewed by the `$ docker trust inspect` command

```sh
$ docker trust inspect --pretty docker.io/openbanking/conformance-suite:v1.0.0

Signatures for docker.io/openbanking/conformance-suite:v1.0.0

SIGNED TAG          DIGEST                                                             SIGNERS
v1.0.0              3d2e482b82608d153a374df3357c0291589a61cc194ec4a9ca2381073a17f58e   jeff

List of signers and their keys for docker.io/openbanking/conformance-suite:v1.0.0

SIGNER              KEYS
jeff                8ae710e3ba82

Administrative keys for docker.io/openbanking/conformance-suite:v1.0.0

  Repository Key:	10b5e94c916a0977471cc08fa56c1a5679819b2005ba6a257aa78ce76d3a1e27
  Root Key:	84ca6e4416416d78c4597e754f38517bea95ab427e5f95871f90d460573071fc
```