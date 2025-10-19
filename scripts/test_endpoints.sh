#!/usr/bin/env bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVER_PORT=${SERVER_PORT:-8080}
BASE_URL="http://localhost:$SERVER_PORT"

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Test result function
test_result() {
    local test_name=$1
    local expected=$2
    local actual=$3
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$expected" = "$actual" ]; then
        echo -e "${GREEN}✓${NC} $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}✗${NC} $test_name"
        echo -e "  Expected: $expected"
        echo -e "  Actual: $actual"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# Print section header
section() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# Main tests
section "Testing Endpoints"

# Test 1: Ping endpoint
echo -n "1. GET /ping... "
RESPONSE=$(curl -s -w "\n%{http_code}" $BASE_URL/ping)
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if test_result "Ping returns 200" "200" "$HTTP_CODE"; then
    if echo "$BODY" | grep -q "pong"; then
        echo -e "   ${GREEN}✓${NC} Response contains 'pong'"
    else
        echo -e "   ${RED}✗${NC} Response missing 'pong': $BODY"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
fi

# Test 2: Login with correct credentials
section "Authentication Tests"

echo -n "2. POST /login (valid credentials)... "
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"password"}')

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

if test_result "Login returns 200" "200" "$HTTP_CODE"; then
    TOKEN=$(echo "$BODY" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$TOKEN" ]; then
        echo -e "   ${GREEN}✓${NC} Token received: ${TOKEN:0:50}..."
    else
        echo -e "   ${RED}✗${NC} No token in response: $BODY"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
fi

# Test 3: Login with incorrect credentials
echo -n "3. POST /login (invalid password)... "
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"wrongpassword"}')

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

if test_result "Invalid login returns 401" "401" "$HTTP_CODE"; then
    if echo "$BODY" | grep -q "error"; then
        echo -e "   ${GREEN}✓${NC} Error message present"
    fi
fi

# Test 4: Access protected endpoint without token
section "Authorization Tests"

echo -n "4. GET /users (no auth)... "
RESPONSE=$(curl -s -w "\n%{http_code}" $BASE_URL/users)
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

test_result "No auth returns 401" "401" "$HTTP_CODE"

# Test 5: Access protected endpoint with invalid token
echo -n "5. GET /users (invalid token)... "
RESPONSE=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer invalid_token" $BASE_URL/users)
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

test_result "Invalid token returns 401" "401" "$HTTP_CODE"

# Test 6-10: Protected endpoints with valid token
if [ -n "$TOKEN" ]; then
    section "Protected Endpoints (with valid token)"
    
    # Test 6: Get users
    echo -n "6. GET /users... "
    RESPONSE=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $TOKEN" $BASE_URL/users)
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')
    
    if test_result "Get users returns 200" "200" "$HTTP_CODE"; then
        if echo "$BODY" | grep -q '"data"'; then
            echo -e "   ${GREEN}✓${NC} Response contains data array"
            
            # Count users
            USER_COUNT=$(echo "$BODY" | grep -o '"username"' | wc -l | tr -d ' ')
            echo -e "   ${BLUE}ℹ${NC} Found $USER_COUNT users"
        fi
    fi
    
    # Test 7: Get approvals
    echo -n "7. GET /approvals... "
    RESPONSE=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $TOKEN" $BASE_URL/approvals)
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')
    
    if test_result "Get approvals returns 200" "200" "$HTTP_CODE"; then
        if echo "$BODY" | grep -q '"data"'; then
            echo -e "   ${GREEN}✓${NC} Response contains data array"
            
            # Count approvals
            APPROVAL_COUNT=$(echo "$BODY" | grep -o '"id"' | wc -l | tr -d ' ')
            echo -e "   ${BLUE}ℹ${NC} Found $APPROVAL_COUNT approvals"
        fi
    fi
    
    # Test 8: Create user (testing POST with data)
    echo -n "8. POST /users (create user)... "
    NEW_USERNAME="test_user_$(date +%s)"
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/users \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$NEW_USERNAME\",\"name\":\"Test User\",\"password\":\"test123\",\"role\":\"staff\"}")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')
    
    if test_result "Create user returns 201" "201" "$HTTP_CODE"; then
        echo -e "   ${GREEN}✓${NC} User created: $NEW_USERNAME"
    fi
    
    # Test 9: Update user
    echo -n "9. PUT /users (update user)... "
    RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT $BASE_URL/users \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"id":1,"name":"Admin Updated"}')
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    
    test_result "Update user returns 200" "200" "$HTTP_CODE"
    
    # Test 10: Create approval
    echo -n "10. POST /approvals (create approval)... "
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/approvals \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    
    test_result "Create approval returns 200" "200" "$HTTP_CODE"
else
    echo -e "${RED}Skipping protected endpoint tests (no token available)${NC}"
fi

# Test summary
section "Test Summary"

echo -e "Total tests:  ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed:       ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed:       ${RED}$FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed! 🎉${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed 😞${NC}"
    exit 1
fi
