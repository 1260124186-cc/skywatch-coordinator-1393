# 修复前故障复现（Docker）

## 项目与标准命令

当前平台为 `linux/arm64`。使用 `docker build -f benzhi.Dockerfile -t go-skywatch-coordinator__011-bug:20260824 .` 构建镜像；镜像内标准编译命令为 `go build ./...`。

## 环境构建与编译

镜像构建成功，容器内 `go build ./...` 通过。

## 故障触发步骤

依次为两个独立活动中的观测站生成计划；第一个站点的计划完成后，再为另一个活动的站点生成计划。

```bash
go test -count=1 -run '^TestStationPlansReleaseCompletedLease$' ./internal/service
```

## 实际错误输出

```text
--- FAIL: TestStationPlansReleaseCompletedLease (0.00s)
    plan_lease_test.go:29: independent second plan should be available: station planning lease is active
FAIL
FAIL	github.com/1260124186-cc/skywatch-coordinator/internal/service	1.377s
FAIL
```

## 期望行为

已完成的站点计划不应阻止其他独立活动继续生成计划；第二次请求应正常可用。
