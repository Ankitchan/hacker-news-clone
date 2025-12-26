#!/bin/bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6ImJvYiIsImVtYWlsIjoiYm9iQGV4YW1wbGUuY29tIiwiaXNzIjoiaGFja2VybmV3cy1jbG9uZSIsImV4cCI6MTc2Njk5OTkyMywibmJmIjoxNzY2NzQwNzIzLCJpYXQiOjE3NjY3NDA3MjN9.ab-sI6OE2xPFdDouwPK9cDHWxpewQwvBiFH7TgxCbLU"

for i in {2..30}; do
  curl -s -X POST http://localhost:8080/api/posts \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"title\":\"Article $i: Latest Tech News and Updates\",\"url\":\"https://news.example.com/article-$i\"}" > /dev/null
  echo "Created post $i"
done

echo ""
echo "Total posts in database:"
curl -s "http://localhost:8080/api/posts" | grep -o '"total_count":[0-9]*' | cut -d: -f2
