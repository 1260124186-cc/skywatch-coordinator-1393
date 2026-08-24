# 修复前故障复现（Docker）

## 项目与标准命令

当前平台为 `linux/arm64`。使用 `docker build -f benzhi.Dockerfile -t go-skywatch-coordinator__013-bug:20260824 .` 构建镜像；镜像内标准编译命令为 `go build ./...`。

## 环境构建与编译

镜像构建成功，容器内 `go build ./...` 通过。

## 故障触发步骤

先取消一次活动观测查询，再对另一个独立活动发起新的查询。

```bash
go test -count=1 -run '^TestCanceledSearchDoesNotPoisonNextCampaignQuery$' ./internal/service
```

## 实际错误输出

```text
--- FAIL: TestCanceledSearchDoesNotPoisonNextCampaignQuery (0.00s)
    query_scope_test.go:22: fresh campaign query should not inherit cancellation: context canceled
FAIL
FAIL	github.com/1260124186-cc/skywatch-coordinator/internal/service	1.173s
FAIL
```

## 期望行为

已取消的查询只应影响自身；后续针对另一活动的新查询应使用有效上下文并正常返回观测结果。
