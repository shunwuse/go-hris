#!/usr/bin/env bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVER_PORT=${SERVER_PORT:-8080}
SERVER_BINARY="./server"
SERVER_LOG="./logs/go-hris-server.log"
PING_TIMEOUT=10
PING_INTERVAL=0.5

# Cleanup function
cleanup() {
    echo -e "${YELLOW}Cleaning up...${NC}"

    # Kill server by PID if available
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        echo "Stopping server (PID: $SERVER_PID)..."
        kill -TERM "$SERVER_PID" 2>/dev/null || true
        sleep 1

        # Force kill if still running
        if kill -0 "$SERVER_PID" 2>/dev/null; then
            echo "Force killing server..."
            kill -KILL "$SERVER_PID" 2>/dev/null || true
        fi
    fi

    # Fallback: kill by port
    if lsof -ti:$SERVER_PORT >/dev/null 2>&1; then
        echo "Killing process on port $SERVER_PORT..."
        lsof -ti:$SERVER_PORT | xargs kill -TERM 2>/dev/null || true
        sleep 1
        lsof -ti:$SERVER_PORT | xargs kill -KILL 2>/dev/null || true
    fi

    # Remove server binary
    if [ -f "$SERVER_BINARY" ]; then
        echo "Removing server binary..."
        rm -f "$SERVER_BINARY"
    fi

    # Remove log file
    if [ -f "$SERVER_LOG" ]; then
        rm -f "$SERVER_LOG"
    fi

    echo -e "${GREEN}Cleanup completed${NC}"
}

# Set trap to cleanup on exit
trap cleanup EXIT INT TERM

# Build server
echo -e "${YELLOW}Building server...${NC}"
go build -o "$SERVER_BINARY" ./cmd/server || {
    echo -e "${RED}Failed to build server${NC}"
    exit 1
}
echo -e "${GREEN}Server built successfully${NC}"

# Start server in background
echo -e "${YELLOW}Starting server on port $SERVER_PORT...${NC}"
"$SERVER_BINARY" > "$SERVER_LOG" 2>&1 &
SERVER_PID=$!

echo "Server PID: $SERVER_PID"
echo "Server log: $SERVER_LOG"

# Wait for server to be ready
echo -e "${YELLOW}Waiting for server to be ready...${NC}"
ELAPSED=0
MAX_ATTEMPTS=$((PING_TIMEOUT * 2))  # Convert to attempts (0.5s interval)
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    if curl -s -f http://localhost:$SERVER_PORT/ping >/dev/null 2>&1; then
        echo -e "${GREEN}Server is ready!${NC}"
        break
    fi

    # Check if server process is still running
    if ! kill -0 "$SERVER_PID" 2>/dev/null; then
        echo -e "${RED}Server process died unexpectedly${NC}"
        echo "Last 20 lines of server log:"
        tail -n 20 "$SERVER_LOG"
        exit 1
    fi

    sleep 0.5
    ATTEMPT=$((ATTEMPT + 1))
    echo -n "."
done
echo ""

# Check if server is ready after timeout
if ! curl -s -f http://localhost:$SERVER_PORT/ping >/dev/null 2>&1; then
    echo -e "${RED}Server failed to start within ${PING_TIMEOUT}s${NC}"
    echo "Server log:"
    cat "$SERVER_LOG"
    exit 1
fi

# Run tests
echo -e "${YELLOW}Running tests...${NC}"
echo ""

# Execute test script if it exists
if [ -f "./scripts/test_endpoints.sh" ]; then
    bash ./scripts/test_endpoints.sh
    TEST_EXIT_CODE=$?
else
    echo -e "${YELLOW}No test_endpoints.sh found, running basic tests...${NC}"

    # Basic ping test
    echo -n "Testing /ping... "
    if curl -s http://localhost:$SERVER_PORT/ping | grep -q "pong"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        TEST_EXIT_CODE=1
    fi

    # Login test
    echo -n "Testing /login... "
    TOKEN=$(curl -s -X POST http://localhost:$SERVER_PORT/login \
        -H "Content-Type: application/json" \
        -d '{"username":"admin","password":"password"}' | \
        grep -o '"token":"[^"]*"' | cut -d'"' -f4)

    if [ -n "$TOKEN" ]; then
        echo -e "${GREEN}✓${NC}"

        # Test /users with token
        echo -n "Testing /users... "
        if curl -s -H "Authorization: Bearer $TOKEN" \
            http://localhost:$SERVER_PORT/users | grep -q '"data"'; then
            echo -e "${GREEN}✓${NC}"
        else
            echo -e "${RED}✗${NC}"
            TEST_EXIT_CODE=1
        fi

        # Test /approvals with token
        echo -n "Testing /approvals... "
        if curl -s -H "Authorization: Bearer $TOKEN" \
            http://localhost:$SERVER_PORT/approvals | grep -q '"data"'; then
            echo -e "${GREEN}✓${NC}"
        else
            echo -e "${RED}✗${NC}"
            TEST_EXIT_CODE=1
        fi
    else
        echo -e "${RED}✗${NC}"
        TEST_EXIT_CODE=1
    fi
fi

echo ""
if [ ${TEST_EXIT_CODE:-0} -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
else
    echo -e "${RED}Some tests failed${NC}"
fi

exit ${TEST_EXIT_CODE:-0}
