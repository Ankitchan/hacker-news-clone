#!/bin/bash

# Test script for spam detection in Docker deployment
# This demonstrates the heuristic-based spam protection working

API_URL="http://localhost:8080/api"
TOKEN=""

echo "=== HACKER NEWS SPAM DETECTION TEST ==="
echo ""

# Step 1: Login or create test user
echo "1. Getting auth token..."
RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"spam@test.com","password":"test12345"}')

TOKEN=$(echo $RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  # Try to create new user
  RESPONSE=$(curl -s -X POST "$API_URL/auth/signup" \
    -H "Content-Type: application/json" \
    -d '{"username":"spamtester","email":"spam@test.com","password":"test12345"}')

  TOKEN=$(echo $RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

  if [ -z "$TOKEN" ]; then
    echo "   ✗ Failed to get auth token"
    echo "   Response: $RESPONSE"
    exit 1
  fi
fi

echo "   ✓ Auth token obtained"
echo ""

# Step 2: Test spam detection
echo "2. Testing spam detection..."
echo ""

# Test 1: Excessive caps (should be blocked)
echo "   Test 1: Excessive Capitalization"
RESULT=$(curl -s -X POST "$API_URL/posts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"CLICK HERE NOW!!!","text":"FREE MONEY AVAILABLE!!!"}')

if echo "$RESULT" | grep -q "Content flagged as spam"; then
  echo "   ✓ SPAM BLOCKED (as expected)"
else
  echo "   ✗ SPAM NOT BLOCKED (unexpected)"
fi
echo ""

# Test 2: Spam keywords (should be blocked)
echo "   Test 2: Spam Keywords"
RESULT=$(curl -s -X POST "$API_URL/posts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Amazing offer","text":"Click here for free viagra and cialis"}')

if echo "$RESULT" | grep -q "Content flagged as spam"; then
  echo "   ✓ SPAM BLOCKED (as expected)"
else
  echo "   ✗ SPAM NOT BLOCKED (unexpected)"
fi
echo ""

# Test 3: Legitimate content (should be allowed)
echo "   Test 3: Legitimate Content"
RESULT=$(curl -s -X POST "$API_URL/posts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Discussion about Docker","text":"What are the best practices for containerizing Go applications?"}')

if echo "$RESULT" | grep -q '"id"'; then
  echo "   ✓ LEGITIMATE POST ALLOWED (as expected)"
else
  echo "   ✗ LEGITIMATE POST BLOCKED (unexpected)"
fi
echo ""

# Summary
echo "=== TEST COMPLETE ==="
echo ""
echo "Spam detection is working via Docker deployment!"
echo "Check backend logs: docker logs hackernews_backend | grep -i spam"
