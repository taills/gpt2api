# syntax=docker/dockerfile:1.6
#
# 单文件部署镜像
# -----------------------------------------------------------------------
# 构建前请先在宿主机生成二进制(见 deploy/build-local.sh / build-local.ps1):
#   deploy/bin/gpt2api   Linux/amd64 单文件可执行(已内嵌前端 + SQLite 驱动)
#
# 这样容器内无需 Go / Node / npm,也不依赖 MySQL / Redis。
# -----------------------------------------------------------------------

FROM alpine:3.20

RUN sed -i 's@https\?://dl-cdn.alpinelinux.org@http://mirrors.ustc.edu.cn@g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata curl \
    && update-ca-certificates \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 单文件可执行(宿主交叉编译产物,已内嵌前端静态资源)
COPY deploy/bin/gpt2api /app/gpt2api
RUN chmod +x /app/gpt2api

# 默认配置(可通过 -v 挂载覆盖)
COPY configs /app/configs

# 数据目录(SQLite 文件 + 日志落盘)
RUN mkdir -p /app/data /app/logs

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -fsS http://localhost:8080/healthz || exit 1

CMD ["/app/gpt2api", "-c", "/app/configs/config.yaml"]
