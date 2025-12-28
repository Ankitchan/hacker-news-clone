# Deploy to Render.com (100% Free)

This guide will help you deploy your Hacker News clone to Render.com for **zero cost**.

## Prerequisites

- GitHub account
- Render.com account (sign up at https://render.com)
- Your code pushed to a GitHub repository

## Step 1: Push Code to GitHub

If you haven't already pushed your code:

```bash
# Make sure you're on the main branch or your preferred branch
git checkout main  # or master

# Push to GitHub
git push origin main
```

If you need to create a new repository:

```bash
# On GitHub, create a new repository (don't initialize it)
# Then:
git remote add origin https://github.com/YOUR_USERNAME/hacker-news-clone.git
git branch -M main
git push -u origin main
```

## Step 2: Sign Up for Render.com

1. Go to https://render.com
2. Click "Get Started"
3. Sign up with GitHub (recommended for easy integration)
4. Authorize Render to access your repositories

## Step 3: Create a Blueprint (Automated Deployment)

### Option A: Using render.yaml (Recommended)

1. **Connect Repository:**
   - In Render dashboard, click "New" → "Blueprint"
   - Select your `hacker-news-clone` repository
   - Render will automatically detect `render.yaml`
   - Click "Apply"

2. **Set Environment Variables:**
   After blueprint is created, go to each service and add:

   **For hackernews-backend:**
   - `CORS_ALLOWED_ORIGINS` = `https://YOUR-FRONTEND-URL.onrender.com`
   - `HF_API_KEY` = `your_huggingface_api_key` (optional, for AI spam detection)

   **For hackernews-frontend:**
   - This will automatically use the backend URL

3. **Wait for Deployment:**
   - Database: ~2 minutes
   - Backend: ~3-5 minutes (Docker build)
   - Frontend: ~2-3 minutes (npm build)

### Option B: Manual Setup (If blueprint doesn't work)

#### 3.1 Create PostgreSQL Database

1. Click "New" → "PostgreSQL"
2. **Name:** `hackernews-db`
3. **Database:** `hackernews_clone`
4. **User:** `postgres`
5. **Region:** Oregon (or closest to you)
6. **Plan:** Free
7. Click "Create Database"
8. **Important:** Copy the "Internal Database URL" - you'll need this!

#### 3.2 Create Backend Service

1. Click "New" → "Web Service"
2. Connect your GitHub repository
3. **Configuration:**
   - **Name:** `hackernews-backend`
   - **Region:** Oregon (same as database)
   - **Branch:** main
   - **Root Directory:** `backend`
   - **Environment:** Docker
   - **Dockerfile Path:** `./Dockerfile`
   - **Plan:** Free

4. **Environment Variables:**
   ```
   PORT=8080
   ENV=production
   DB_HOST=<from Internal Database URL>
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=<from Internal Database URL>
   DB_NAME=hackernews_clone
   DB_SSL_MODE=require
   JWT_SECRET=<generate a random 32+ character string>
   JWT_EXPIRATION_HOURS=72
   CORS_ALLOWED_ORIGINS=https://YOUR-FRONTEND-URL.onrender.com
   SPAM_DETECTION_ENABLED=true
   HF_API_KEY=<your_huggingface_key_optional>
   HF_SPAM_MODEL=mrm8488/bert-tiny-finetuned-sms-spam-detection
   ```

5. Click "Create Web Service"

#### 3.3 Create Frontend Service

1. Click "New" → "Static Site"
2. Connect your GitHub repository
3. **Configuration:**
   - **Name:** `hackernews-frontend`
   - **Region:** Oregon
   - **Branch:** main
   - **Root Directory:** `frontend`
   - **Build Command:** `npm ci && npm run build`
   - **Publish Directory:** `dist`
   - **Plan:** Free

4. **Environment Variables:**
   ```
   VITE_API_BASE_URL=https://YOUR-BACKEND-URL.onrender.com/api
   ```
   (Replace with your actual backend URL from step 3.2)

5. Click "Create Static Site"

## Step 4: Update CORS After Deployment

Once frontend is deployed, you'll get a URL like:
`https://hackernews-frontend-XXXX.onrender.com`

1. Go to your backend service settings
2. Update `CORS_ALLOWED_ORIGINS` to include your frontend URL:
   ```
   https://hackernews-frontend-XXXX.onrender.com
   ```
3. Click "Save Changes"
4. Backend will automatically redeploy

## Step 5: Test Your Deployment

1. Visit your frontend URL: `https://hackernews-frontend-XXXX.onrender.com`
2. Try creating an account
3. Create a post
4. Test spam detection by trying to post with "viagra" or "cialis"

## Important Notes

### Free Tier Limitations

- **Services sleep after 15 minutes of inactivity**
  - First request after sleep will take ~30-60 seconds to wake up
  - Subsequent requests will be fast

- **512 MB RAM per service**
  - Sufficient for this application

- **PostgreSQL free tier:**
  - 1 GB storage
  - Expires after 90 days (you'll need to recreate it)
  - Or upgrade to paid ($7/month for persistent database)

### Keep Your App Alive (Optional)

If you want to prevent sleeping, you can:

1. Use a free service like UptimeRobot to ping your site every 14 minutes
2. Upgrade to paid plan ($7/month for backend, keeps it running 24/7)

### Database Backups

Free tier databases are deleted after 90 days. To backup:

```bash
# Export from Render dashboard
# Or use pg_dump with the external connection string
```

## Troubleshooting

### Backend won't start
- Check logs in Render dashboard
- Verify all environment variables are set
- Ensure database is running and accessible

### Frontend can't connect to backend
- Verify `VITE_API_BASE_URL` is correct
- Check `CORS_ALLOWED_ORIGINS` includes your frontend URL
- Check backend logs for CORS errors

### Database connection errors
- Ensure `DB_SSL_MODE=require` for Render PostgreSQL
- Check database is running (not suspended)
- Verify database credentials in environment variables

### Spam detection not working
- Check `HF_API_KEY` is set (optional, heuristics work without it)
- View backend logs for spam detection errors
- Heuristic checks still work even without API key

## Custom Domain (Optional)

Render supports custom domains on free tier:

1. Go to your frontend service → "Settings" → "Custom Domain"
2. Add your domain (e.g., `hackernews.yourdomain.com`)
3. Follow DNS instructions to point your domain to Render
4. Update backend CORS to include your custom domain

## Estimated Deployment Time

- Total setup time: 15-20 minutes
- Initial deployment: 10-15 minutes
- Subsequent deployments (on git push): 5-10 minutes

## Cost Summary

- **Everything:** $0/month
- **Optional upgrades:**
  - Persistent PostgreSQL: $7/month
  - Always-on backend: $7/month
  - Professional plan (faster builds): $19/month

## Next Steps

After deployment:
1. Test all features thoroughly
2. Set up monitoring (Render provides basic monitoring)
3. Consider enabling custom domain
4. Share your deployed app!

## Support

- Render Documentation: https://render.com/docs
- Render Community: https://community.render.com
- Check backend logs in Render dashboard for issues

---

**Your app will be live at:**
- Frontend: `https://hackernews-frontend-XXXX.onrender.com`
- Backend: `https://hackernews-backend-XXXX.onrender.com`

Enjoy your free deployment! 🎉
