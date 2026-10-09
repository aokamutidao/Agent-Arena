#!/usr/bin/env bash
# Course application only. It never kills processes it did not start.
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
if [ -f "$ROOT_DIR/.env.course" ]; then
  set -a
  source "$ROOT_DIR/.env.course"
  set +a
fi
JWT_SECRET="${JWT_SECRET:-}"
if [ "${#JWT_SECRET}" -lt 32 ]; then
  echo 'Set JWT_SECRET (at least 32 random characters) in .env.course. See .env.course.example.' >&2
  exit 1
fi
COURSE_API_PORT="${PORT:-8080}"
COURSE_UI_PORT="${FRONTEND_PORT:-3000}"
export PORT="$COURSE_API_PORT"
export NEXT_PUBLIC_API_URL="${NEXT_PUBLIC_API_URL-http://localhost:$COURSE_API_PORT}"
export NEXT_TELEMETRY_DISABLED=1
COURSE_API_PID=''
COURSE_UI_PID=''
cleanup() {
  [ -z "$COURSE_UI_PID" ] || kill "$COURSE_UI_PID" 2>/dev/null || true
  [ -z "$COURSE_API_PID" ] || kill "$COURSE_API_PID" 2>/dev/null || true
}
trap cleanup EXIT
trap 'exit 130' INT TERM
cd "$ROOT_DIR/backend"
mkdir -p tmp
go build -o tmp/course-server ./cmd/course-server
./tmp/course-server &
COURSE_API_PID=$!
COURSE_READY=false
for attempt in $(seq 1 30); do
  if ! kill -0 "$COURSE_API_PID" 2>/dev/null; then
    echo 'Course backend failed to start; inspect the error above.' >&2
    exit 1
  fi
  if curl --fail --silent "http://127.0.0.1:$COURSE_API_PORT/health" >/dev/null; then
    COURSE_READY=true
    break
  fi
  sleep 1
done
if [ "$COURSE_READY" != true ]; then
  echo 'Course backend health check timed out.' >&2
  exit 1
fi
cd "$ROOT_DIR/frontend"
npm run dev -- --port "$COURSE_UI_PORT" &
COURSE_UI_PID=$!
echo "Course UI: http://localhost:$COURSE_UI_PORT"
echo "Readiness: http://localhost:$COURSE_API_PORT/ready"
echo 'An unconfigured chain remains read-only; deploy the new course contracts to enable transactions.'
wait "$COURSE_UI_PID"
