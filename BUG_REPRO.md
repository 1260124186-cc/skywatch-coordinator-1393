# 修复前故障复现（Docker）

## 项目与标准命令

当前平台为 `linux/arm64`。使用 `docker build -f benzhi.Dockerfile -t go-skywatch-coordinator__014-bug:20260824 .` 构建镜像；镜像内标准编译命令为 `go build ./...`。

## 环境构建与编译

镜像构建成功，容器内 `go build ./...` 通过。

## 故障触发步骤

先读取北区活动的汇总，再读取南区活动的汇总；两个活动的站点和班次数据互不相同。

```bash
go test -count=1 -run '^TestCampaignSummaryDoesNotReusePreviousScope$' ./internal/service
```

## 实际错误输出

```text
--- FAIL: TestCampaignSummaryDoesNotReusePreviousScope (0.00s)
    summary_scope_test.go:33: second summary = domain.CampaignSummary{Campaign:domain.Campaign{ID:"cmp-000001", Name:"North Summary"}, StationCount:1, ShiftCount:1, OpenShiftCount:1}, want campaign "cmp-000004" with no shifts
FAIL
FAIL	github.com/1260124186-cc/skywatch-coordinator/internal/service	0.927s
FAIL
```

## 期望行为

活动汇总必须按当前活动独立计算；第二次读取应返回南区活动及其自身的统计数据，而非沿用北区活动的汇总。
