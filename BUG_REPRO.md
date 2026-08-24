# Skywatch Coordinator 审核派发阻塞复现

## 项目与标准命令

Skywatch Coordinator 是用于协调野外天文观测活动的 Go 服务。

## 环境构建与编译

```bash
go build ./...
```

## 故障触发步骤

```bash
go test -count=1 -run '^TestReviewsDoNotShareRelayAcrossCampaigns$' ./internal/service
```

## 实际错误输出

第二个独立活动的审核返回 `review relay is busy`。

## 期望行为

一个活动完成审核后不应占用另一活动的审核派发状态，两个活动都应能独立保存审核结果。
