#!/usr/bin/env bash
# gpt2api 容器启动入口。
#
# 单文件部署:无外部数据库依赖,SQLite 在首次启动时自动初始化。
# 此脚本仅做环境检查后直接 exec 主进程。

set -euo pipefail

log() { echo "[entrypoint] $*"; }

# 确保数据目录存在且可写
DATA_DIR="${GPT2API_DATA_DIR:-/app/data}"
mkdir -p "$DATA_DIR"
log "data dir: $DATA_DIR"

log "starting: $*"
exec "$@"
