# Hacker News Clone API Documentation

Base URL: `http://localhost:8080/api`

## Table of Contents
- [Authentication](#authentication)
- [Posts](#posts)
- [Comments](#comments)
- [Votes](#votes)
- [Error Responses](#error-responses)

---

## Authentication

### Sign Up
Create a new user account.

**Endpoint:** `POST /api/auth/signup`

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:** `201 Created`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Validation:**
- Username: minimum 3 characters
- Password: minimum 8 characters, maximum 72 characters
- Email and username must be unique

---

### Login
Authenticate an existing user.

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### Get Profile
Get authenticated user's profile.

**Endpoint:** `GET /api/auth/profile`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

### Refresh Token
Get a new JWT token.

**Endpoint:** `POST /api/auth/refresh`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## Posts

### Create Post
Create a new post (authenticated).

**Endpoint:** `POST /api/posts`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body (URL post):**
```json
{
  "title": "Interesting Article",
  "url": "https://example.com/article"
}
```

**Request Body (Text post):**
```json
{
  "title": "Discussion Topic",
  "text": "Let's discuss this interesting topic..."
}
```

**Request Body (Both URL and Text):**
```json
{
  "title": "Article with Commentary",
  "url": "https://example.com/article",
  "text": "Here's my take on this..."
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "title": "Interesting Article",
  "url": "https://example.com/article",
  "text": null,
  "user_id": 1,
  "username": "johndoe",
  "points": 0,
  "comment_count": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### Get All Posts
Retrieve all posts with pagination.

**Endpoint:** `GET /api/posts`

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Response:** `200 OK`
```json
{
  "posts": [
    {
      "id": 1,
      "title": "Interesting Article",
      "url": "https://example.com/article",
      "user_id": 1,
      "username": "johndoe",
      "points": 10,
      "comment_count": 5,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total_count": 100,
  "page": 1,
  "page_size": 20
}
```

---

### Get Post by ID
Retrieve a single post.

**Endpoint:** `GET /api/posts/{id}`

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Interesting Article",
  "url": "https://example.com/article",
  "user_id": 1,
  "username": "johndoe",
  "points": 10,
  "comment_count": 5,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### Get New Posts
Retrieve posts sorted by newest first.

**Endpoint:** `GET /api/posts/new`

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Response:** Same as Get All Posts

---

### Get Top Posts
Retrieve posts sorted by points.

**Endpoint:** `GET /api/posts/top`

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Response:** Same as Get All Posts

---

### Get Best Posts
Retrieve posts sorted by best algorithm (points/age).

**Endpoint:** `GET /api/posts/best`

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Response:** Same as Get All Posts

---

### Search Posts
Search posts by title or text content.

**Endpoint:** `GET /api/posts/search`

**Query Parameters:**
- `q` (required): Search query
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Example:** `GET /api/posts/search?q=javascript&page=1`

**Response:** Same as Get All Posts

---

### Update Post
Update an existing post (owner only).

**Endpoint:** `PUT /api/posts/{id}` or `PATCH /api/posts/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "title": "Updated Title",
  "text": "Updated content..."
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Updated Title",
  "text": "Updated content...",
  "user_id": 1,
  "username": "johndoe",
  "points": 10,
  "comment_count": 5,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T01:00:00Z"
}
```

---

### Delete Post
Delete a post (owner only).

**Endpoint:** `DELETE /api/posts/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Post deleted successfully"
}
```

---

## Comments

### Create Comment
Create a new comment (authenticated).

**Endpoint:** `POST /api/comments`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body (Top-level comment):**
```json
{
  "post_id": 1,
  "text": "Great article!"
}
```

**Request Body (Reply to comment):**
```json
{
  "post_id": 1,
  "parent_id": 5,
  "text": "I agree with your point."
}
```

**Response:** `201 Created`
```json
{
  "id": 10,
  "post_id": 1,
  "user_id": 1,
  "username": "johndoe",
  "parent_id": null,
  "text": "Great article!",
  "points": 0,
  "depth": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### Get Comment by ID
Retrieve a single comment.

**Endpoint:** `GET /api/comments/{id}`

**Response:** `200 OK`
```json
{
  "id": 10,
  "post_id": 1,
  "user_id": 1,
  "username": "johndoe",
  "parent_id": null,
  "text": "Great article!",
  "points": 5,
  "depth": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### Get Comments by Post
Retrieve all comments for a post.

**Endpoint:** `GET /api/posts/{post_id}/comments`

**Query Parameters:**
- `threaded` (optional, default: "true"): "true" for nested structure, "false" for flat list

**Response (Threaded):** `200 OK`
```json
{
  "comments": [
    {
      "id": 10,
      "post_id": 1,
      "user_id": 1,
      "username": "johndoe",
      "parent_id": null,
      "text": "Great article!",
      "points": 5,
      "depth": 0,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "children": [
        {
          "id": 11,
          "post_id": 1,
          "user_id": 2,
          "username": "janedoe",
          "parent_id": 10,
          "text": "I agree!",
          "points": 2,
          "depth": 1,
          "created_at": "2024-01-01T00:05:00Z",
          "updated_at": "2024-01-01T00:05:00Z",
          "children": []
        }
      ]
    }
  ]
}
```

---

### Get Comment Replies
Get direct replies to a comment.

**Endpoint:** `GET /api/comments/{id}/replies`

**Response:** `200 OK`
```json
[
  {
    "id": 11,
    "post_id": 1,
    "user_id": 2,
    "username": "janedoe",
    "parent_id": 10,
    "text": "I agree!",
    "points": 2,
    "depth": 1,
    "created_at": "2024-01-01T00:05:00Z",
    "updated_at": "2024-01-01T00:05:00Z"
  }
]
```

---

### Update Comment
Update a comment (owner only).

**Endpoint:** `PUT /api/comments/{id}` or `PATCH /api/comments/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "text": "Updated comment text"
}
```

**Response:** `200 OK`
```json
{
  "id": 10,
  "post_id": 1,
  "user_id": 1,
  "username": "johndoe",
  "parent_id": null,
  "text": "Updated comment text",
  "points": 5,
  "depth": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T01:00:00Z"
}
```

---

### Delete Comment
Delete a comment (owner only).

**Endpoint:** `DELETE /api/comments/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Comment deleted successfully"
}
```

---

## Votes

All vote endpoints require authentication.

### Vote on Post
Upvote or downvote a post.

**Endpoint:** `POST /api/votes/posts/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "vote_type": 1
}
```

- `vote_type`: `1` for upvote, `-1` for downvote

**Response:** `200 OK`
```json
{
  "success": true,
  "points": 15
}
```

---

### Vote on Comment
Upvote or downvote a comment.

**Endpoint:** `POST /api/votes/comments/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "vote_type": 1
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "points": 8
}
```

---

### Remove Vote from Post
Remove your vote from a post.

**Endpoint:** `DELETE /api/votes/posts/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "success": true,
  "points": 14
}
```

---

### Remove Vote from Comment
Remove your vote from a comment.

**Endpoint:** `DELETE /api/votes/comments/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "success": true,
  "points": 7
}
```

---

### Get User's Vote on Post
Check if the current user has voted on a post.

**Endpoint:** `GET /api/votes/posts/{id}/user`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (has voted):** `200 OK`
```json
{
  "id": 100,
  "user_id": 1,
  "post_id": 1,
  "vote_type": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Response (no vote):** `200 OK`
```json
{
  "vote": null
}
```

---

### Get User's Vote on Comment
Check if the current user has voted on a comment.

**Endpoint:** `GET /api/votes/comments/{id}/user`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** Same as Get User's Vote on Post

---

## Error Responses

All error responses follow this format:

**Response:** `4xx` or `5xx`
```json
{
  "error": "Bad Request",
  "message": "Detailed error message"
}
```

### Common Error Codes

- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Insufficient permissions (e.g., trying to edit someone else's post)
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., email already exists)
- `500 Internal Server Error`: Server error

---

## Health Check

### Check API Health
Verify the API is running.

**Endpoint:** `GET /api/health`

**Response:** `200 OK`
```json
{
  "status": "ok",
  "message": "Hacker News Clone API is running"
}
```

---

## Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Tokens are obtained from:
- `/api/auth/signup` - When creating a new account
- `/api/auth/login` - When logging in
- `/api/auth/refresh` - When refreshing an existing token

Tokens expire after 72 hours by default (configurable via `JWT_EXPIRATION_HOURS` environment variable).
