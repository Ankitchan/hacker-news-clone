import requests
import time

# Login
response = requests.post('http://localhost:8080/api/auth/login', json={
    'username': 'alice',
    'password': 'password123'
})

if response.status_code != 200:
    print(f"Login failed: {response.text}")
    exit(1)

token = response.json()['token']
print(f"Logged in successfully")

# Create posts
for i in range(2, 31):
    response = requests.post(
        'http://localhost:8080/api/posts',
        headers={'Authorization': f'Bearer {token}'},
        json={
            'title': f'Interesting Article {i}: How Technology is Changing Our World',
            'url': f'https://example.com/article-{i}'
        }
    )
    if response.status_code == 201:
        print(f"Created post {i}")
    else:
        print(f"Failed to create post {i}: {response.text}")
    time.sleep(0.05)

print("Done!")
