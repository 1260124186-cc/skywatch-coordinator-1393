# skywatch-coordinator__009 Docker 交付说明

## 项目概览
- Skywatch Coordinator is a small HTTP service used by field-science teams to coordinate night-sky observation campaigns. A coordinator opens a campaign, stations operate observation
- Go module: `github.com/1260124186-cc/skywatch-coordinator`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/skywatch
```

## Docker 构建

```bash
./build_benzhi_docker.sh skywatch-coordinator__009-benzhi linux/amd64
docker run --rm -it skywatch-coordinator__009-benzhi bash
```

## 环境

- 基础镜像: `golang:1.26.2`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
- 源码中检测到的服务端口: `18087`
