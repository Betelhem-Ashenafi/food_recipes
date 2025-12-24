# 🚀 Quick Vercel Deployment Guide

Deploy your Food Recipes App to Vercel in minutes!

## 📋 Prerequisites

1. **GitHub Account** - Your code should be on GitHub
2. **Vercel Account** - Sign up at [vercel.com](https://vercel.com) (free)
3. **Backend deployed** - You'll need your backend API running somewhere (see options below)

---

## 🎯 Step-by-Step Deployment

### Step 1: Prepare Your Code

Make sure your code is pushed to GitHub:
```bash
git add .
git commit -m "Prepare for Vercel deployment"
git push origin main
```

### Step 2: Deploy Frontend to Vercel

1. **Go to [vercel.com](https://vercel.com)** and sign in
2. **Click "Add New Project"**
3. **Import your GitHub repository**
   - Select your `food-recipes-app` repository
   - Click "Import"

4. **Configure Project Settings:**
   - **Framework Preset**: Nuxt.js (auto-detected)
   - **Root Directory**: `frontend/nuxt-app`
   - **Build Command**: `npm run build`
   - **Output Directory**: `.output/public`
   - **Install Command**: `npm install`

5. **Add Environment Variables:**
   Click "Environment Variables" and add:
   ```
   NUXT_PUBLIC_API_URL=https://your-backend-url.vercel.app
   NUXT_PUBLIC_HASURA_URL=https://your-hasura-url.com/v1/graphql
   ```

6. **Click "Deploy"**

7. **Wait for deployment** (2-3 minutes)

8. **Get your link!** 
   - Vercel will give you a URL like: `https://your-app.vercel.app`
   - Share this with your boss! 🎉

---

## 🔧 Backend Deployment Options

Since Vercel is best for frontend, you have 3 options for the backend:

### Option 1: Railway (Easiest - Recommended)

1. Go to [railway.app](https://railway.app)
2. Sign up with GitHub
3. New Project → Deploy from GitHub
4. Select your repo
5. Add service:
   - **Backend**: Root directory: `backend`
   - **PostgreSQL**: Add database
   - **Hasura**: Use Docker image `hasura/graphql-engine:v2.40.0`
6. Set environment variables
7. Get backend URL and update frontend env vars in Vercel

### Option 2: Render (Free Tier Available)

1. Go to [render.com](https://render.com)
2. New → Web Service
3. Connect GitHub repo
4. Settings:
   - Root Directory: `backend`
   - Build: `go build -o main .`
   - Start: `./main`
5. Add PostgreSQL database
6. Get URL and update Vercel env vars

### Option 3: Keep Backend Local (For Testing Only)

If you just need a quick demo:
- Keep backend running locally
- Use [ngrok](https://ngrok.com) to expose it:
  ```bash
  ngrok http 8081
  ```
- Use the ngrok URL in Vercel environment variables

---

## ⚙️ Environment Variables Setup

### In Vercel (Frontend):

Go to Project Settings → Environment Variables:

```bash
# Backend API URL
NUXT_PUBLIC_API_URL=https://your-backend.railway.app

# Hasura GraphQL URL  
NUXT_PUBLIC_HASURA_URL=https://your-hasura.railway.app/v1/graphql

# Optional: Hasura Admin Secret (for local dev only)
NUXT_PUBLIC_HASURA_ADMIN_SECRET=myhasurasecret
```

### In Backend Service (Railway/Render):

```bash
# Database
DATABASE_URL=postgres://user:pass@host:5432/dbname

# JWT Secret
JWT_SECRET=your-super-secret-key-here

# Port (usually auto-set)
PORT=8081
```

---

## 🔄 Update Frontend Code for Production

Update API URLs in your frontend code to use environment variables:

The code should already use `process.env.NUXT_PUBLIC_API_URL` or hardcoded `localhost:8081`. 

**Quick Fix**: Update all `http://localhost:8081` to use environment variable:

```javascript
const API_URL = process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8081';
```

---

## 📝 Quick Deployment Checklist

- [ ] Code pushed to GitHub
- [ ] Vercel account created
- [ ] Frontend deployed to Vercel
- [ ] Backend deployed (Railway/Render)
- [ ] Environment variables set in Vercel
- [ ] Environment variables set in backend service
- [ ] Database migrations run
- [ ] Test the deployed app
- [ ] Share link with boss! 🎉

---

## 🎁 Your Deployment Link

After deployment, Vercel will give you:
- **Production URL**: `https://your-app.vercel.app`
- **Preview URLs**: For each commit/PR

**Share this with your boss!** 👔

---

## 🐛 Troubleshooting

### Frontend not loading?
- Check environment variables in Vercel
- Check build logs in Vercel dashboard
- Make sure `NUXT_PUBLIC_*` variables are set

### API calls failing?
- Check backend is running
- Check CORS settings in backend
- Verify API URL in environment variables

### Database errors?
- Make sure database is running
- Check connection string
- Run migrations if needed

---

## 💡 Pro Tips

1. **Custom Domain**: Add your own domain in Vercel settings
2. **Preview Deployments**: Every PR gets its own preview URL
3. **Analytics**: Enable Vercel Analytics to track usage
4. **Speed**: Vercel CDN makes your app super fast globally

---

**Good luck with your deployment! 🚀**

