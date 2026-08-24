# Skywatch Coordinator 分析扫描阻塞复现

## 项目与标准命令

Skywatch Coordinator 协调野外观测活动。

## 环境构建与编译

```bash
go build ./...
```

## 故障触发步骤

```bash
go test -count=1 -run '^TestAnalyticsScansStayIndependentAcrossCampaigns$' ./internal/service
```

## 实际错误输出

第二个活动返回 `analytics scan is busy`。

## 期望行为

不同活动的统计请求应独立完成。
