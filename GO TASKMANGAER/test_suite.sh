#!/bin/bash
set -e

PORT=8085
BASE_URL="http://localhost:${PORT}/api/v1"

echo "=== Starting Integration Test Suite ==="

# Build binary
echo "[1] Building TaskManager binary..."
go build -o taskmanager .

# Run server in background
echo "[2] Starting TaskManager server on port ${PORT}..."
./taskmanager > server_test.log 2>&1 &
SERVER_PID=$!

cleanup() {
    echo "Stopping server (PID: $SERVER_PID)..."
    kill $SERVER_PID || true
}
trap cleanup EXIT

# Wait for server to start
sleep 2

# Test 1: Register User 1
echo "[3] Testing User 1 Registration..."
REG_RESP=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"Alice","email":"alice@example.com","password":"secretpassword"}')
HTTP_CODE=$(echo "$REG_RESP" | tail -n1)
BODY=$(echo "$REG_RESP" | sed '$d')
if [ "$HTTP_CODE" -ne 201 ]; then
    echo "FAIL: Expected status 201 on registration, got $HTTP_CODE. Response: $BODY"
    exit 1
fi
echo "PASS: Registered User 1 (Alice)."

# Test 2: Register Duplicate Email
echo "[4] Testing Duplicate Email Registration..."
DUP_RESP=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"Alice Duplicate","email":"alice@example.com","password":"secretpassword"}')
HTTP_CODE=$(echo "$DUP_RESP" | tail -n1)
if [ "$HTTP_CODE" -ne 400 ]; then
    echo "FAIL: Expected status 400 on duplicate registration, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Rejected duplicate email with 400."

# Test 3: Login User 1
echo "[5] Testing User 1 Login..."
LOGIN_RESP=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"alice@example.com","password":"secretpassword"}')
TOKEN1=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*' | grep -o '[^"]*$')
if [ -z "$TOKEN1" ]; then
    echo "FAIL: Could not obtain token from login. Response: $LOGIN_RESP"
    exit 1
fi
echo "PASS: Obtained JWT token for Alice."

# Test 4: Access protected route without token
echo "[6] Testing Protected Route Without Token..."
UNAUTH_RESP=$(curl -s -w "\n%{http_code}" "${BASE_URL}/tasks")
HTTP_CODE=$(echo "$UNAUTH_RESP" | tail -n1)
if [ "$HTTP_CODE" -ne 401 ]; then
    echo "FAIL: Expected status 401 without token, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Unauthorized access blocked with 401."

# Test 5: Create Task as User 1
echo "[7] Testing Task Creation..."
CREATE_RESP=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/tasks" \
    -H "Authorization: Bearer $TOKEN1" \
    -H "Content-Type: application/json" \
    -d '{"title":"Complete Milestone 2 & 3","description":"Build Postgres & Auth layer"}')
HTTP_CODE=$(echo "$CREATE_RESP" | tail -n1)
BODY=$(echo "$CREATE_RESP" | sed '$d')
if [ "$HTTP_CODE" -ne 201 ]; then
    echo "FAIL: Expected status 201 on task creation, got $HTTP_CODE. Response: $BODY"
    exit 1
fi
TASK_ID=$(echo "$BODY" | grep -o '"ID":[0-9]*' | head -n1 | cut -d: -f2)
echo "PASS: Created Task ID $TASK_ID for Alice."

# Test 6: Get User 1 Tasks
echo "[8] Testing Get Tasks..."
GET_TASKS_RESP=$(curl -s -w "\n%{http_code}" "${BASE_URL}/tasks?status=pending" \
    -H "Authorization: Bearer $TOKEN1")
HTTP_CODE=$(echo "$GET_TASKS_RESP" | tail -n1)
if [ "$HTTP_CODE" -ne 200 ]; then
    echo "FAIL: Expected status 200 on GetTasks, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Retrieved tasks list successfully."

# Test 7: Update Task with Invalid Status (Validation test)
echo "[9] Testing Status Validation on Update..."
INVALID_STATUS_RESP=$(curl -s -w "\n%{http_code}" -X PUT "${BASE_URL}/tasks/${TASK_ID}" \
    -H "Authorization: Bearer $TOKEN1" \
    -H "Content-Type: application/json" \
    -d '{"title":"Test","status":"invalid_status"}')
HTTP_CODE=$(echo "$INVALID_STATUS_RESP" | tail -n1)
if [ "$HTTP_CODE" -ne 422 ]; then
    echo "FAIL: Expected status 422 on invalid status update, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Invalid status rejected with 422 Unprocessable Entity."

# Test 8: Valid Update Task (PUT)
echo "[10] Testing Valid Task Update (PUT)..."
UPDATE_RESP=$(curl -s -w "\n%{http_code}" -X PUT "${BASE_URL}/tasks/${TASK_ID}" \
    -H "Authorization: Bearer $TOKEN1" \
    -H "Content-Type: application/json" \
    -d '{"title":"Milestone 2 & 3 Completed","description":"Updated description","status":"in_progress"}')
HTTP_CODE=$(echo "$UPDATE_RESP" | tail -n1)
if [ "$HTTP_CODE" -ne 200 ]; then
    echo "FAIL: Expected status 200 on task update, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Task updated successfully with PUT."

# Test 9: Patch Task (PATCH)
echo "[11] Testing Patch Task (PATCH)..."
PATCH_RESP=$(curl -s -w "\n%{http_code}" -X PATCH "${BASE_URL}/tasks/${TASK_ID}" \
    -H "Authorization: Bearer $TOKEN1" \
    -H "Content-Type: application/json" \
    -d '{"status":"done"}')
HTTP_CODE=$(echo "$PATCH_RESP" | tail -n1)
if [ "$HTTP_CODE" -ne 200 ]; then
    echo "FAIL: Expected status 200 on task patch, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Task status patched to 'done'."

# Test 10: User 2 Ownership Isolation
echo "[12] Testing Task Ownership Isolation with User 2..."
curl -s -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"Bob","email":"bob@example.com","password":"secretpassword"}' > /dev/null

LOGIN_BOB=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"bob@example.com","password":"secretpassword"}')
TOKEN2=$(echo "$LOGIN_BOB" | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

# Bob tries to view Alice's task
BOB_GET=$(curl -s -w "\n%{http_code}" "${BASE_URL}/tasks/${TASK_ID}" \
    -H "Authorization: Bearer $TOKEN2")
HTTP_CODE=$(echo "$BOB_GET" | tail -n1)
if [ "$HTTP_CODE" -ne 403 ]; then
    echo "FAIL: Expected status 403 when Bob attempts to access Alice's task, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Ownership isolation enforced (403 Forbidden)."

# Test 11: User 1 Deletes Task
echo "[13] Testing Task Deletion..."
DEL_RESP=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}/tasks/${TASK_ID}" \
    -H "Authorization: Bearer $TOKEN1")
HTTP_CODE=$(echo "$DEL_RESP" | tail -n1)
if [ "$HTTP_CODE" -ne 204 ]; then
    echo "FAIL: Expected status 204 on task deletion, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Task deleted with 204 No Content."

# Test 12: Refresh Token
echo "[14] Testing Refresh Token..."
REFRESH_RESP=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/auth/refresh" \
    -H "Authorization: Bearer $TOKEN1")
HTTP_CODE=$(echo "$REFRESH_RESP" | tail -n1)
BODY=$(echo "$REFRESH_RESP" | sed '$d')
NEW_TOKEN=$(echo "$BODY" | grep -o '"token":"[^"]*' | grep -o '[^"]*$')
if [ "$HTTP_CODE" -ne 200 ] || [ -z "$NEW_TOKEN" ]; then
    echo "FAIL: Expected status 200 and new token on refresh, got $HTTP_CODE."
    exit 1
fi
echo "PASS: Refreshed token successfully."

echo "=== ALL INTEGRATION TESTS PASSED SUCCESSFULLY! ==="
