# common-bandwidth-auto-switch

![wtfpl](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-1.png)
[中文文档](readme.zh.md)

## why

Alibaba Cloud’s fancy 95-tier shared bandwidth plan guarantees a minimum of 200Mbps, with a minimum spend of 200*0.2=40Mbps. After some serious number crunching (and a bit of guesswork), I figured that using it somewhere between 40~50Mbps is the sweet spot.

Thanks to the 2020 Chinese New Year + Wuhan’s infamous coronavirus party crasher, website traffic plummeted faster than my motivation on a Monday morning. So, I shoved all the IPs into the shared bandwidth pool.

On 2020-03-05 at 13:20 (Thursday, because why not), Baidu decided to crawl our site like a caffeine-fueled spider, and bandwidth instantly shot up to 200Mbps.

I originally planned to yank the EIP out of the shared bandwidth. But that afternoon, freshly woken from a nap, operating at a solid 10% brainpower, I misread the screen and accidentally removed the SLB bound to the EIP, causing the site to go *poof* and become unreachable.

![img](/img/guo.jpg)

After licking my wounds, I vowed to develop an automatic program to manage and adjust the shared bandwidth EIP pool so that it can:

1. During low traffic periods, shove EIPs into the shared bandwidth to save money;
1. During bandwidth rush hour, yank EIPs out of the shared bandwidth to boost bandwidth utilization.

All in all, to squeeze out every last drop of cost optimization.

![img](/img/b.jpg)

## feature

When current shared bandwidth > desired threshold (maxBandwidth), automatically kick out the high-bandwidth EIPs from the shared bandwidth pool.

When current shared bandwidth < desired threshold ((maxBandwidth+minBandwidth)/2), automatically invite high-bandwidth EIPs back into the shared bandwidth pool.

The core algorithm is `dynamic programming`. Fancy, right?

Once bandwidth reaches equilibrium, you’ll realize this program is about as useful as a chocolate teapot — it barely triggers any scaling at all.

![](/img/fly.jpg)

## usage

### RAM Permissions

APIs used:

```bash
vpc
DescribeEipMonitorData
DescribeEipAddresses

cms
DescribeCommonBandwidthPackages
AddCommonBandwidthPackageIp
RemoveCommonBandwidthPackageIp

cms
DescribeMetricList
```

Because you’re messing with VPC and shared bandwidth — aka **high-risk operations** — make sure your RAM permissions are on point.

If you’re feeling lazy about permissions, just slap on these:

1. AliyunVPCFullAccess
1. AliyunEIPFullAccess
1. AliyunCloudMonitorFullAccess
1. AliyunCommonBandwidthPackageFullAccess

### Configuration File

Copy `config-example.yaml` to `config.yaml` and tweak it to your heart’s content.

```bash
cp config-example.yaml config.yaml
vi config.yaml
```

### docker/docker-compose

```bash
# docker
docker run -it \
-v /root/common-bandwidth-auto-switch/config.yaml:/app/config/config.yaml \
zeusro/common-bandwidth-auto-switch:latest
# docker-compose
docker-compose up
```

### kubernetes

Here’s an example of deploying with Kubernetes using a `CronJob` + `Configmap` combo, because why not?

```bash
kubectl create configmap cbwp-config --from-file=config.yaml=config.yaml

kubectl apply -f deploy/kubernetes/cronjob.yaml
```

Configuration loading order: tries `config.yaml` first; if missing, falls back to `config-example.yaml`.

## rm -rf /

```bash
kubectl delete cronjob common-bandwidth-auto-switch
kubectl delete configmap cbwp-config
```

## todo(NEVER DO)

1. Make it blocking, run periodically, and add health checks.
1. Run multiple instances.

Because who doesn’t love living on the edge?
