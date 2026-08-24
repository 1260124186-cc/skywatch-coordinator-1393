# 修复前故障复现（Docker）

## 环境构建与编译

使用当前 `linux/arm64` 环境构建，标准编译命令为 `go build ./...`。

## 故障触发步骤

依次执行两个互不相关活动的操作：北区巡检或重叠检查结束后，南区独立活动的相同检查因为前一场遗留的活动范围而收到冲突。

```bash
go test -count=1 -run '^TestInspectionReleasesCampaignScope$' ./internal/service
```

## 实际错误输出

第二个独立活动的请求返回冲突/忙碌错误。

## 期望行为

独立活动的巡检完成后必须清除上一场活动的检查范围，南区巡检不应被北区残留状态拒绝。
