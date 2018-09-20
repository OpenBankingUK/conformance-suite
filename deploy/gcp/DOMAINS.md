# `DOMAINS.md`

**NB:** In the future this can be replaced by using this Kubernetes addon, https://github.com/kubernetes-incubator/external-dns. For now, we will just create the records manually. See this script [./domains-create.sh](./domains-create.sh).

We are using [Google Cloud DNS](https://cloud.google.com/dns/) to create the records, e.g., here are the already created records for the <https://compliance-suite-server.ga> domain:

```sh
$ gcloud dns record-sets list --zone=compliance-suite-server-ga-zone
NAME                             TYPE   TTL    DATA
compliance-suite-server.ga.      A      60     35.230.158.145
compliance-suite-server.ga.      NS     21600  ns-cloud-e1.googledomains.com.,ns-cloud-e2.googledomains.com.,ns-cloud-e3.googledomains.com.,ns-cloud-e4.googledomains.com.
compliance-suite-server.ga.      SOA    21600  ns-cloud-e1.googledomains.com. cloud-dns-hostmaster.google.com. 47 21600 3600 259200 300
www.compliance-suite-server.ga.  CNAME  60     compliance-suite-server.ga.
```

This maps:

* compliance-suite-server.ga to the LoadBalancer IP, 35.230.158.145.
* www.compliance-suite-server.ga maps to previous `A` record, compliance-suite-server.ga.
* Setups the NameServers, `NS`, to point to the Google Cloud DNS ones.

Trial run:

```sh
ping compliance-suite-server.ga
PING compliance-suite-server.ga (35.230.158.145): 56 data bytes
...
```
