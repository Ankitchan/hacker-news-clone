#!/bin/bash

set -e

echo "=== Testing Notification System ==="
echo ""

# Login as Alice (post owner)
echo "1. Logging in as Alice..."
ALICE_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password123"}' | jq -r '.token')
echo "Alice logged in successfully"
echo ""

# Login as Bob (commenter)
echo "2. Logging in as Bob..."
BOB_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"bob@example.com","password":"password123"}' | jq -r '.token')
echo "Bob logged in successfully"
echo ""

# Alice creates a post
echo "3. Alice creates a new post..."
POST_RESPONSE=$(curl -s -X POST http://localhost:8080/api/posts \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{"title":"Test Post for Notifications","url":"https://example.com/test"}')
POST_ID=$(echo $POST_RESPONSE | jq -r '.id')
echo "Post created with ID: $POST_ID"
echo ""

# Bob comments on Alice's post (should trigger notification)
echo "4. Bob comments on Alice's post..."
COMMENT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/posts/$POST_ID/comments \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{"text":"Great post, Alice!"}')
COMMENT_ID=$(echo $COMMENT_RESPONSE | jq -r '.id')
echo "Comment created with ID: $COMMENT_ID"
echo ""

# Alice checks her notifications
echo "5. Alice checks her notifications..."
NOTIFICATIONS=$(curl -s -X GET http://localhost:8080/api/notifications \
  -H "Authorization: Bearer $ALICE_TOKEN")
echo "Notifications:"
echo $NOTIFICATIONS | jq '.'
echo ""

# Get unread count
echo "6. Alice checks unread notification count..."
UNREAD_COUNT=$(curl -s -X GET http://localhost:8080/api/notifications/unread/count \
  -H "Authorization: Bearer $ALICE_TOKEN")
echo "Unread count:"
echo $UNREAD_COUNT | jq '.'
echo ""

# Alice replies to Bob's comment
echo "7. Alice replies to Bob's comment..."
REPLY_RESPONSE=$(curl -s -X POST http://localhost:8080/api/comments \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d "{\"post_id\":$POST_ID,\"parent_id\":$COMMENT_ID,\"text\":\"Thanks, Bob!\"}")
REPLY_ID=$(echo $REPLY_RESPONSE | jq -r '.id')
echo "Reply created with ID: $REPLY_ID"
echo ""

# Bob checks his notifications (should have notification about Alice's reply)
echo "8. Bob checks his notifications..."
BOB_NOTIFICATIONS=$(curl -s -X GET http://localhost:8080/api/notifications \
  -H "Authorization: Bearer $BOB_TOKEN")
echo "Bob's notifications:"
echo $BOB_NOTIFICATIONS | jq '.'
echo ""

# Mark notifications as read
echo "9. Alice marks her notifications as read..."
NOTIFICATION_ID=$(echo $NOTIFICATIONS | jq -r '.notifications[0].id')
MARK_READ=$(curl -s -X POST http://localhost:8080/api/notifications/mark-read \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d "{\"notification_ids\":[$NOTIFICATION_ID]}")
echo "Mark as read response:"
echo $MARK_READ | jq '.'
echo ""

# Verify notification is marked as read
echo "10. Verify Alice's notification is marked as read..."
UPDATED_NOTIFICATIONS=$(curl -s -X GET http://localhost:8080/api/notifications \
  -H "Authorization: Bearer $ALICE_TOKEN")
echo "Updated notifications:"
echo $UPDATED_NOTIFICATIONS | jq '.notifications[0].is_read'
echo ""

echo "=== Test Complete ==="
