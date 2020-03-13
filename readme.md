# common-bandwidth-auto-switch

## why

阿里云增强型95付费的共享带宽,保底200Mbps,最低消费200*0.2=40Mbps,经过我计算,我觉得在40~50Mbps时使用比较合理.

由于2020年春节+武汉冠状肺炎的严重影响,网站流量急剧减少.所以我把所有的IP都纳入了共享带宽.

2020-03-05 13:20星期四,百度在爬我们站点,带宽瞬间达到200Mbps.

我原本想把EIP脱离共享带宽,结果却错误地移除了EIP绑定的SLB,造成了网站无法访问.

![](/img/guo.jpg)

最终我决定,开发一个自动管理调整共享带宽EIP池的玩意.

1. 流量低峰时把EIP纳入共享带宽,节约流量费用;
1. 带宽高峰时让EIP移出共享带宽,提高带宽利用率.

最终最大程度地优化费用支出.

![](/img/b.jpg)

## feature

当前共享带宽 > 期望值时,自动把高带宽的EIP移除出共享带宽

当前共享带宽 < 期望值时,自动把高带宽的EIP添加入共享带宽

## usage

例子使用了 Kubernetes 的 `CronJob` + `secret` 的方式部署

```bash
# edit your config
kubectl create secret generic cbwp-config --from-file=config.yaml=config-example.yaml
# or
kubectl create secret generic cbwp-config --from-file=config.yaml=config.yaml
kubectl apply -f deploy/kubernetes/cronjob.yaml
```

我在`.dockerignore`里面留了一手,没有忽略掉 `config-example.yaml` 文件,也可以配置这个文件,然后直接把配置打包到容器里面,这样volume都不需要配置

配置的加载顺序为:先读取 `config.yaml` , `config.yaml`不存在再读取 `config-example.yaml`文件.

只是说出于安全,不太这么建议这么做.

## warning

用到的接口:

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

由于需要操作VPC和共享带宽，这类都属于**高危操作**，RAM授权记得弄好。

## todo(todo means NEVER DO)

1. 阻塞,周期性运行,再加上健康检查
1. 多实例运行