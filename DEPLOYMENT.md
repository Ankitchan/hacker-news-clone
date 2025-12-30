# Deployment Guide

This guide covers deploying the Hacker News Clone application using Docker and Docker Compose.

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- 2GB RAM minimum
- 10GB disk space

## Quick Start with Docker Compose

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd hacker_news_clone
```

### 2. Configure Environment Variables

Copy the example environment file:

```bash
cp .env.docker .env
```

Edit `.env` and update the following **critical** values:

```bash
# IMPORTANT: Change these for production!
JWT_SECRET=your_super_secret_jwt_key_at_least_32_characters_change_this
DB_PASSWORD=postgres_secure_password_change_me
```

### 3. Optional: Enable Spam Detection

To enable AI-powered spam detection:

1. Get a free API key from [Hugging Face](https://huggingface.co/settings/tokens)
2. Update `.env`:

```bash
SPAM_DETECTION_ENABLED=true
HF_API_KEY=your_huggingface_api_key_here
```

**Note:** The free tier includes monthly credits. The spam detector also has built-in heuristics that work even without the API.

### 4. Start the Application

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Backend API on port 8080
- Frontend on port 80

### 5. Access the Application

Open your browser and navigate to:
- **Frontend:** http://localhost
- **Backend API:** http://localhost:8080

### 6. Create Initial Users

The database starts empty. Register your first user through the frontend at:
http://localhost/signup

## Deployment

### Currently deployed on: Render.com

1. Create a `render.yaml`:

```yaml
services:
  - type: web
    name: hackernews-backend
    env: docker
    dockerfilePath: ./backend/Dockerfile
    envVars:
      - key: DATABASE_URL
        fromDatabase:
          name: hackernews-db
          property: connectionString
      - key: JWT_SECRET
        generateValue: true
      - key: PORT
        value: 8080

  - type: web
    name: hackernews-frontend
    env: docker
    dockerfilePath: ./frontend/Dockerfile

databases:
  - name: hackernews-db
    databaseName: hackernews
    user: postgres
```

2. Connect your repository to Render.com
3. Deploy with one click


## Manual Docker Commands

### Build Images

```bash
# Build backend
docker build -t hackernews-backend ./backend

# Build frontend
docker build -t hackernews-frontend ./frontend
```

### Run Individual Containers

```bash
# Run PostgreSQL
docker run -d \
  --name hackernews_db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=hackernews_clone \
  -p 5432:5432 \
  postgres:16-alpine

# Run backend
docker run -d \
  --name hackernews_backend \
  --link hackernews_db:db \
  -e DB_HOST=db \
  -e JWT_SECRET=your_secret \
  -p 8080:8080 \
  hackernews-backend

# Run frontend
docker run -d \
  --name hackernews_frontend \
  -p 80:80 \
  hackernews-frontend
```

## Monitoring and Logs

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f db
```

### Service Status

```bash
docker-compose ps
```

### Resource Usage

```bash
docker stats
```

## Maintenance

### Backup Database

```bash
docker exec hackernews_db pg_dump -U postgres hackernews_clone > backup.sql
```

### Restore Database

```bash
cat backup.sql | docker exec -i hackernews_db psql -U postgres hackernews_clone
```

### Update Application

```bash
git pull
docker-compose down
docker-compose build
docker-compose up -d
```

### Clean Up

```bash
# Stop and remove containers
docker-compose down

# Remove volumes (WARNING: This deletes all data!)
docker-compose down -v

# Remove images
docker-compose down --rmi all
```

## Troubleshooting

### Backend Won't Start

Check logs:
```bash
docker-compose logs backend
```

Common issues:
- Database not ready: Wait for health check
- Missing env variables: Check `.env` file
- Port already in use: Change PORT in `.env`

### Frontend 404 Errors

Nginx configuration issue. Check:
```bash
docker-compose logs frontend
```

Rebuild frontend:
```bash
docker-compose up -d --build frontend
```

### Database Connection Failed

Ensure database is healthy:
```bash
docker-compose ps db
```

Check connectivity:
```bash
docker exec hackernews_backend ping db
```

### Spam Detection Not Working

1. Verify `SPAM_DETECTION_ENABLED=true` in `.env`
2. Check `HF_API_KEY` is set correctly
3. View logs for API errors:
   ```bash
   docker-compose logs backend | grep -i spam
   ```

## Performance Tuning

### Scale Services

```bash
# Run multiple backend instances
docker-compose up -d --scale backend=3
```

### Database Performance

Add to `docker-compose.yml` under `db` service:

```yaml
command:
  - "postgres"
  - "-c"
  - "shared_buffers=256MB"
  - "-c"
  - "max_connections=200"
```

### Minimize Spam Detection Costs

- Set `SPAM_DETECTION_ENABLED=false` to disable API calls
- Built-in heuristics still provide basic protection
- Or use Hugging Face free tier (monthly credits)

## Production Recommendations

1. **Use a reverse proxy** (Nginx/Caddy) for HTTPS
2. **Set up monitoring** (Prometheus + Grafana)
3. **Configure log rotation**
4. **Set up automated backups**
5. **Use managed PostgreSQL** (AWS RDS, DigitalOcean Managed DB)
6. **Enable spam detection** with Hugging Face API key
7. **Set up CI/CD** for automated deployments

## Support

For issues:
- Check logs: `docker-compose logs`
- Restart services: `docker-compose restart`
- Full reset: `docker-compose down && docker-compose up -d`

For more help, open an issue on GitHub.
