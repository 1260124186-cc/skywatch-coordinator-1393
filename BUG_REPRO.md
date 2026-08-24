# 私有复现说明

建立两个时间和站点均独立的观测活动，各自提交一个不同标签。先读取第一场活动的标签，再读取第二场活动的标签。Bug 环境把第一场的标签带入第二场响应。

复现命令：

`go test -count=1 -run '^TestObservationLabelsRemainScopedToTheirCampaign$' ./internal/service`
