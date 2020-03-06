# common-bandwidth-auto-switch

共享带宽EIP自调整,流量低峰时把EIP纳入共享带宽,节约流量费用;带宽高峰时让EIP移出共享带宽.

最大程度地优化费用支出.

## feature

当前共享带宽 > 期望值时,自动把高带宽的EIP移除出共享带宽

当前共享带宽 < 1/2期望值时,自动把高带宽的EIP添加入共享带宽

## warning

由于需要操作VPC和共享带宽，这类都属于**高危操作**，RAM授权记得弄好。

## example

## usage

## todo(todo means NEVER DO)

1. 阻塞,周期性运行

