#!/usr/bin/env bash
# gpt2api 本地预构建脚本 — 生成单文件可执行 (Linux/amd64)
#
# 用法:
#   bash deploy/build-local.sh            # 常规构建
#   bash deploy/build-local.sh --no-web   # 跳过前端构建(仅重编后端)
#
# 产物:
#   deploy/bin/gpt2api    linux/amd64 单文件可执行(已内嵌前端)

set -euo pipefail

SKIP_WEB=0
for arg in "$@"; do
    case "$arg" in
        --no-web) SKIP_WEB=1 ;;
    esac
done

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
EMBED_DIR="internal/server/web"

echo "[build-local] repo = $ROOT"

# ---- step1: 前端 ----
if [ "$SKIP_WEB" = "0" ]; then
    echo "[build-local] step1 = npm run build (web)"
    pushd web >/dev/null
    if [ ! -d node_modules ]; then
        npm install --no-audit --no-fund --loglevel=error
    fi
    npm run build
    popd >/dev/null

    echo "[build-local] step1 = copy dist → ${EMBED_DIR}/"
    rm -rf "$EMBED_DIR"
    mkdir -p "$EMBED_DIR"
    cp -r web/dist/. "$EMBED_DIR/"
else
    echo "[build-local] step1 = skipped (--no-web)"
fi

# ---- step2: 后端(embed 已内嵌前端) ----
echo "[build-local] step2 = cross-build gpt2api (linux/amd64)"
mkdir -p deploy/bin
# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
CGO_ENABLED=0 \
    go build -ldflags "-s -w" -o deploy/bin/gpt2api ./cmd/server

echo "[build-local] done."
ls -lh deploy/bin/gpt2api
