#!/bin/bash
# Agent Arena 一键启动脚本
# 同时启动后端 (8080) 和前端 (3000)

set -e

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 确保 go/node/npm 在 PATH 中
export PATH="/usr/local/go/bin:/usr/local/bin:$PATH"

# 颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}🏟️  Agent Arena 启动中...${NC}"

# 清理旧进程
lsof -ti:8080 2>/dev/null | xargs kill -9 2>/dev/null || true
lsof -ti:3000 2>/dev/null | xargs kill -9 2>/dev/null || true
sleep 1

# 启动后端
echo -e "${BLUE}[1/2] 启动后端 :8080${NC}"
cd "$PROJECT_DIR/backend"
go run ./cmd/server/ &
BACKEND_PID=$!

# 等待后端就绪
for i in $(seq 1 15); do
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "${GREEN}  ✓ 后端就绪${NC}"
        break
    fi
    sleep 1
done

# 启动前端
echo -e "${BLUE}[2/2] 启动前端 :3000${NC}"
cd "$PROJECT_DIR/frontend"
npm run dev &
FRONTEND_PID=$!

# 等待前端就绪
for i in $(seq 1 20); do
    if curl -s -o /dev/null http://localhost:3000 2>/dev/null; then
        echo -e "${GREEN}  ✓ 前端就绪${NC}"
        break
    fi
    sleep 1
done

echo ""
echo -e "${GREEN}🏟️  Agent Arena 已启动！${NC}"
echo -e "  前端: ${BLUE}http://localhost:3000${NC}"
echo -e "  后端: ${BLUE}http://localhost:8080${NC}"
echo -e "  对局: ${BLUE}http://localhost:3000/game/1${NC}"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 捕获退出信号，清理子进程
trap "echo '正在停止...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0" INT TERM

# 等待子进程
wait
