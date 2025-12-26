#!/bin/bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VybmFtZSI6ImFsaWNlIiwiZW1haWwiOiJhbGljZUBleGFtcGxlLmNvbSIsImV4cCI6MTczNTQ3Nzg1MywiaWF0IjoxNzM1MjE4NjUzfQ.QMoWaMGYSZQlb9QOcQiQ4bT8zOeBcPr1dhtA6ckHqbo"

for i in {2..30}; do
  curl -X POST http://localhost:8080/api/posts \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"title\":\"Article $i: Technology News\",\"url\":\"https://example.com/$i\"}" 2>/dev/null
  echo "Post $i created"
  sleep 0.1
done
