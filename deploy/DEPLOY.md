# `DEPLOY.md`
## live
* production: [https://compliance-suite-server.ga](https://compliance-suite-server.ga).

### urls
* <https://compliance-suite-server.ga>

### services
**NB:** none of these are exposed to the external world at present. this is intentional.

| service               | url                                                   | notes                         |
|-----------------------|-------------------------------------------------------|-------------------------------|
| kafka                 | kafka.compliance-suite-server.tk:9092                 | Temporarily disabled for now. |
| zookeeper             | zookeeper.compliance-suite-server.tk:2181             | Temporarily disabled for now. |
| redis                 | redis.compliance-suite-server.tk:6379                 | Temporarily disabled for now. |
| mongo                 | mongo.compliance-suite-server.tk:27017                | Temporarily disabled for now. |
| reference-mock-server | reference-mock-server.compliance-suite-server.tk:8001 | Temporarily disabled for now. |
| tpp-reference-server  | tpp-reference-server.compliance-suite-server.tk:8003  | Temporarily disabled for now. |

## quick start
installs go, kompose. **NB:** we use the `master` branch of `kompose` as the release version tends ot fall behind quiet a bit.

```sh
$ echo -e "\033[92m  ---> installing go \033[0m"
$ brew install go
$ mkdir -p $HOME/go
$ cat >> $HOME/.bash_profile <<-"EOF"
##### go config
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
#####
EOF
$ source $HOME/.bash_profile
$ go version
$ echo -e "\033[92m  ---> installing kompose \033[0m"
$ go get -u github.com/kubernetes/kompose && kompose version
```

### convert from Docker Compose to Kubernetes manifests
```sh
$ make convert
```

### updating the live deploy to reference new images
#### manually
TODO: document how to manually push updates.

### debugging
See [./DEPLOY-DEBUG.md](./DEPLOY-DEBUG.md).
