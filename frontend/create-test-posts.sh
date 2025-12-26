#!/bin/bash

# Get the token from alice's login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"password123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "Token: $TOKEN"

# Create 25 test posts
for i in {2..26}; do
  echo "Creating post $i..."
  curl -s -X POST http://localhost:8080/api/posts \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"title\": \"Interesting Article $i: How Technology is Changing Our World\",
      \"url\": \"https://example.com/article-$i\"
    }" > /dev/null
  sleep 0.1
done

echo "Created 25 test posts!"
echo "Total posts in database:"
curl -s http://localhost:8080/api/posts | python3 -c "import sys, json; data=json.load(sys.stdin); print(f'Total: {data[\"total_count\"]}')"
