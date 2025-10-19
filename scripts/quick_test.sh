#!/usr/bin/env bash

# Quick test script - assumes server is already running
# Usage: ./scripts/quick_test.sh [base_url]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BASE_URL=${1:-"http://localhost:8080"}

echo -e "${YELLOW}Quick API Test${NC}"
echo -e "Testing: $BASE_URL"
echo ""

# Ping test
echo -n "Ping... "
if curl -s -f $BASE_URL/ping | grep -q "pong"; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Login test
echo -n "Login... "
TOKEN=$(curl -s -X POST $BASE_URL/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"password"}' | \
    grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Users test
echo -n "Get users... "
if curl -s -H "Authorization: Bearer $TOKEN" $BASE_URL/users | grep -q '"data"'; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Approvals test
echo -n "Get approvals... "
if curl -s -H "Authorization: Bearer $TOKEN" $BASE_URL/approvals | grep -q '"data"'; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}All quick tests passed!${NC}"
