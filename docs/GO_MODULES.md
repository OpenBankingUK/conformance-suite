# `GO_MODULES.md`

## FAQ

> Just interested in a defined recovery strategy for the situation when a package disappears from github for example. Unlikely be if we define what happens ahead of time, reduces the risk and panic. And from a governance point of view it covers that base. So we need a recovery strategy defined.

**References:**

* https://golang.org/cmd/go/#hdr-Module_proxy_protocol
* https://goproxy.io/
* https://github.com/golang/go/wiki/Modules#are-there-always-on-module-repositories-and-enterprise-proxies
* https://groups.google.com/d/msg/golang-dev/mNedL5rYLCs/OGjRDTmWBgAJ
* https://github.com/thepudds/go-module-knobs/blob/master/README.md: "Note that the go command stores downloaded dependencies in its local cache ($GOPATH/pkg/mod) and its cache layout is the same as the requirements for a proxy, so that cache can be used as the content for a filesystem-based GOPROXY or simple webserver used as a GOPROXY."

I don’t know what to do when a package disappears. Maybe look for an archived version?
