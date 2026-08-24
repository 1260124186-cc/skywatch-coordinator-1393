# 修复前故障复现（Docker）

## 项目与标准命令

当前平台为 linux/arm64。使用 `docker build -f benzhi.Dockerfile -t go-skywatch-coordinator__001-bug:20260824 .` 构建镜像；镜像内标准编译命令为 `go build ./...`。

## 环境构建与编译

镜像构建成功，容器内 `go version` 输出 `go version go1.26.2 linux/arm64`，并且 `go build ./...` 通过。

## 故障触发步骤

在包含验证用例的工作区执行：

```bash
go test -count=1 -run '^TestReleaseCampaignsKeepIndependentDispatchState$' ./internal/service
```

先完成一场活动的发布，再发布另一场已完成审核且时段已关闭的活动。

## 实际错误输出

```text
--- FAIL: TestReleaseCampaignsKeepIndependentDispatchState (0.00s)
    release_signal_test.go:35: independent second release should not inherit dispatch state: release dispatch is busy
FAIL
FAIL	github.com/1260124186-cc/skywatch-coordinator/internal/service	0.544s
FAIL
```

## 期望行为

不同活动的连续发布应各自完成，后一场活动不应继承前一场发布后的忙碌状态。
