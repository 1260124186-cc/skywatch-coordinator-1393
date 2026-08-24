# 修复前故障复现（Docker）

## 环境构建与编译

使用当前 `linux/arm64` 环境构建，标准编译命令为 `go build ./...`。

## 故障触发步骤

依次执行两个互不相关活动的操作：完成北区的站点初始化或班次开启后，南区独立活动的后续初始化操作被运行期门禁拒绝。

```bash
go test -count=1 -run '^TestStationsReleaseCompletedSetupGate$' ./internal/service
```

## 实际错误输出

第二个独立活动的请求返回冲突/忙碌错误。

## 期望行为

独立活动的站点配置完成后必须释放初始化占用，不能阻止下一场活动继续新增站点。
