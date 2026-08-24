# Skywatch Coordinator 班次关闭阻塞复现

## 项目与标准命令

Skywatch Coordinator 协调野外观测活动。

## 环境构建与编译

```bash
go build ./...
```

## 故障触发步骤

```bash
go test -count=1 -run '^TestShiftClosuresStayIndependentAcrossCampaigns$' ./internal/service
```

## 实际错误输出

第二个独立班次关闭时返回 `shift closure is busy`。

## 期望行为

不同活动的班次关闭状态应独立，均可正常记录结束时间。
