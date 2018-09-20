# `DEPLOY-DEBUG.md`
### debugging
#### volumes
This example lists the contents of a volume that stores the [Let’s Encrypt](https://letsencrypt.org/) certificates. Listing contents of `PersistentVolumeClaim`/`PersistentVolume`.

1. SSH into Compute Engine instance that has the `PersistentVolumeClaim`/`PersistentVolume` attached.
2. Find out where the the `gcePersistentDisk`, in this case its `/var/lib/kubelet/plugins/kubernetes.io/gce-pd/mounts/gke-compliance-suite-s-pvc-6d014a5e-891e-11e8-9454-42010a9a01d1`, is mounted:

```
$ sudo df -h
Filesystem      Size  Used Avail Use% Mounted on
/dev/root       1.2G  465M  757M  39% /
devtmpfs        846M     0  846M   0% /dev
tmpfs           848M     0  848M   0% /dev/shm
tmpfs           848M  1.3M  847M   1% /run
tmpfs           848M     0  848M   0% /sys/fs/cgroup
tmpfs           848M     0  848M   0% /tmp
tmpfs           256K     0  256K   0% /mnt/disks
/dev/sda8        12M   28K   12M   1% /usr/share/oem
/dev/sda1        26G  4.4G   22G  18% /mnt/stateful_partition
overlayfs       1.0M  164K  860K  17% /etc
tmpfs           1.0M  128K  896K  13% /var/lib/cloud
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/416e5745-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/metadata-agent-token-6zgcc
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/fb14177174c042e6c51c1c5257232cd9bcd9ea78b13895555651eb880e09f179/merged
shm              64M     0   64M   0% /var/lib/docker/containers/62a70ecc41bf5d43f9a2f7dadef4332d54277c2d1dbdee718aa38df0af738a1d/shm
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/adbd441d570a5974135194de66c43433189185795259042bffb380819ec81812/merged
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/de3afec99770effc604a1153adc1b55d9b31109aa9a10d0ae77f519777abab20/merged
shm              64M     0   64M   0% /var/lib/docker/containers/dd7779ee9dba2931c8f4d17460de75213f6c372adb20d91d3326a1a85ff40e43/shm
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/3e191916-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/kubernetes-dashboard-token-t6k9l
tmpfs           848M     0  848M   0% /var/lib/kubelet/pods/3e191916-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/kubernetes-dashboard-certs
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/c4ce0727453d49429860d72a3ab403a4f2f0a12336a52cebce73e01262634374/merged
shm              64M     0   64M   0% /var/lib/docker/containers/a031e16ff0f25814438bb705c72f836095b261c9b341b80ec3140fee58aaf26a/shm
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/ee4a86e313d88b39e5f0b3095b9c57749ae7e96cb2f2d322400a2b3dc21a85fa/merged
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/63a7b3fa-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/metrics-server-token-bdg4j
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/a2755532e5175db2a61efd6fff07880d0c3188ba67093389126f95546a98ce39/merged
shm              64M     0   64M   0% /var/lib/docker/containers/c710e53be94ffc3e2fb40e04bb00215968488f01aca2fe73f83db83f50a87b59/shm
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/3ae71be0cb42bfca3114a3b2000f98c626228977dc2bc61df38a3dcba8aa6022/merged
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/7cb09fab5f7d6d4728dcf72975471ce9569743523fcd4c005a88f03ee18610c2/merged
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/8ed68ac5-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/fluentd-gcp-token-rm2sn
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/ca2ec922d9ee6de5c9ec207698fce2a64b7fb2d92d7723f40eac22bd61a57ed4/merged
shm              64M     0   64M   0% /var/lib/docker/containers/6f9a6b84df3fa6644a851f12439ec8611b234af0a96d2f584a8b1b94bfb70068/shm
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/9c9d0c5c0c437f34170f4f9be8d7ea840f3d0b6702548ec9bd7e69373e1e6f67/merged
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/8c16aed0d77210b33745d9fc3b1fc067662005edaabc4a75eadb6c8c47e23ddc/merged
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/9e3357cf-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/default-token-xbc7c
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/b979ba5c8fb14e97b203f76a1cbdf200ef1f27d1a643ad9cd32ea95b55654880/merged
shm              64M     0   64M   0% /var/lib/docker/containers/c69c5e0d2116c74a552e11ebf110ac72035ccc99e3618c1d121f52851e9e2204/shm
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/9e8259f1-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/default-token-xbc7c
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/20d5305f38c97c8432531a25acece0780b4c6bac467aaa4bc39bd67fa228ea74/merged
shm              64M     0   64M   0% /var/lib/docker/containers/9daa2d2d8a8668574d8e74a487da8f0c6ca80cf7743a351e389c931613da2a49/shm
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/9ec44f59-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/default-token-xbc7c
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/261cfd294068b8cd01600c88c1c53d772e99eff7064c65e952e1087efe8246ea/merged
shm              64M     0   64M   0% /var/lib/docker/containers/e2abf145a2557d6593f41a4c7dbfb069ab5d7cd8eda2eb721d3026a82de83705/shm
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/9f5dfef8-891e-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/default-token-xbc7c
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/2aaa0305bbc2eaf184b7651e9ff234e33ab448f94ac97c690b6cc868a39e9b47/merged
shm              64M     0   64M   0% /var/lib/docker/containers/3b8257b961f15f4c27ca17d3adab5044f3abb23c3c21d76b417b006565d8d924/shm
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/f10e7b0f58c6622b7b60220f1ee71ae84aeddaf457bb82883658958bfdd19f1d/merged
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/a635755343b95dd0246fec6ddc5d3b440c5f40591e55e5356743437ac18781d6/merged
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/ee4994e88a335910abf55ed89947650c86593b6fb65ac5dec77697f79aeab652/merged
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/0fb888618ebc65457f59d6bbaa81db55152634d16935958ef4e9c7a9cb9049b5/merged
tmpfs           848M   12K  848M   1% /var/lib/kubelet/pods/042f4b0d-8938-11e8-9454-42010a9a01d1/volumes/kubernetes.io~secret/default-token-xbc7c
/dev/sdb        976M  2.6M  907M   1% /var/lib/kubelet/plugins/kubernetes.io/gce-pd/mounts/gke-compliance-suite-s-pvc-6d014a5e-891e-11e8-9454-42010a9a01d1
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/f9ff94bc8ecacfaa4dba2bfb4e70174638ebf1850693a3d04b03bb0f8d3a0f11/merged
shm              64M     0   64M   0% /var/lib/docker/containers/1219acb07db3f11b86c2ce4e069fc31b4a6ec6362375d06cabdcb3c711895059/shm
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/297f676167a7c80554bca40dcb3958ff69be59347ec0e0f58ebe74825453bdc0/merged
overlay          26G  4.4G   22G  18% /var/lib/docker/overlay2/a558e9791aef4ee3b0829c984ed3c86ec10c78836dd7b4f49d13c2a512f418a9/merged
```

3. Print contents of the volume:

```
$ sudo ls -lah /var/lib/kubelet/plugins/kubernetes.io/gce-pd/mounts/gke-compliance-suite-s-pvc-6d014a5e-891e-11e8-9454-42010a9a01d1
total 32K
drwxr-xr-x 5 root root 4.0K Jul 16 17:37 .
drwxr-x--- 3 root root 4.0K Jul 16 20:37 ..
drwx------ 4 root root 4.0K Jul 16 17:35 acme
drwx------ 2 root root  16K Jul 16 17:35 lost+found
drwx------ 2 root root 4.0K Jul 16 17:37 ocsp
```

which is the directory structure [Caddy](https://caddyserver.com/) creates for all the certificates it has downloaded.
