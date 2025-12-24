# 🚀 Step-by-Step Vercel Deployment Guide

Follow these steps **one by one** to deploy your app and get a link for your boss.

---

## 📦 STEP 1: Prepare Your Code

**Do this:**
1. Open terminal/PowerShell in your project folder
2. Run these commands:

```bash
cd D:\food-recipes-app
git add .
git commit -m "Ready for Vercel deployment"
git push origin main
```

**Check:** Make sure your code is pushed to GitHub (go to github.com and verify you see your latest code)

---

## 🌐 STEP 2: Create Vercel Account

**Do this:**
1. Go to: https://vercel.com
2. Click "Sign Up" (top right)
3. Choose "Continue with GitHub"
4. Authorize Vercel to access your GitHub
5. Complete signup

**Check:** You should see the Vercel dashboard

---

## 🎯 STEP 3: Deploy Frontend to Vercel

**Do this:**
1. In Vercel dashboard, click **"Add New..."** button (top right)
2. Click **"Project"**
3. You'll see your GitHub repositories
4. Find **"food-recipes-app"** and click **"Import"**

**Do this next:**
5. On the "Configure Project" page, set these:
   - **Framework Preset**: Should auto-detect "Nuxt.js" ✅
   - **Root Directory**: Click "Edit" and type: `frontend/nuxt-app`
   - **Build Command**: Should be `npm run build` ✅
   - **Output Directory**: Should be `.output/public` ✅
   - **Install Command**: Should be `npm install` ✅

6. **DON'T add environment variables yet** - click **"Deploy"** button

**Wait:** 2-3 minutes for deployment to complete

**Check:** You should see "Building..." then "Ready" with a green checkmark

---

## 🔗 STEP 4: Get Your Frontend Link

**Do this:**
1. After deployment completes, you'll see a success page
2. Look for a URL like: `https://food-recipes-app-xxxxx.vercel.app`
3. **Copy this URL** - this is your frontend link!

**Check:** Click the link - you should see your app (but it won't work fully yet because backend isn't connected)

---

## 🗄️ STEP 5: Deploy Backend to Railway

**Do this:**
1. Go to: https://railway.app
2. Click "Start a New Project"
3. Click "Login with GitHub"
4. Authorize Railway

**Do this next:**
5. Click "New Project"
6. Click "Deploy from GitHub repo"
7. Select your **"food-recipes-app"** repository
8. Click "Deploy Now"

**Wait:** Railway will start deploying

---

## 🐘 STEP 6: Add PostgreSQL Database

**Do this:**
1. In Railway project, click **"+ New"** button
2. Click **"Database"**
3. Click **"Add PostgreSQL"**
4. Wait 30 seconds for database to create

**Check:** You should see a PostgreSQL service in your project

**Do this next:**
5. Click on the PostgreSQL service
6. Go to **"Variables"** tab
7. Find **"DATABASE_URL"** - **COPY THIS VALUE** (you'll need it later)

---

## ⚙️ STEP 7: Add Backend Service

**Do this:**
1. In Railway project, click **"+ New"** button again
2. Click **"GitHub Repo"**
3. Select **"food-recipes-app"**
4. Railway will ask for settings:
   - **Root Directory**: Type `backend`
   - Click **"Deploy"**

**Do this next:**
5. After it starts, click on the backend service
6. Go to **"Variables"** tab
7. Click **"+ New Variable"** and add these:

   **Variable 1:**
   - Name: `DATABASE_URL`
   - Value: Paste the DATABASE_URL you copied from PostgreSQL
   - Click "Add"

   **Variable 2:**
   - Name: `JWT_SECRET`
   - Value: `my-super-secret-jwt-key-123456789`
   - Click "Add"

   **Variable 3:**
   - Name: `PORT`
   - Value: `8081`
   - Click "Add"

8. Railway will auto-redeploy with new variables

**Wait:** 2-3 minutes for backend to deploy

**Do this next:**
9. Go to **"Settings"** tab of backend service
10. Find **"Generate Domain"** button and click it
11. **COPY THE URL** (e.g., `https://backend-production-xxxx.up.railway.app`)

---

## 🔷 STEP 8: Add Hasura Service

**Do this:**
1. In Railway project, click **"+ New"** button
2. Click **"Empty Service"**
3. Click on the new service
4. Go to **"Settings"** tab
5. Scroll to **"Docker"** section
6. Enable **"Use Docker"**
7. In **"Docker Image"** field, type: `hasura/graphql-engine:v2.40.0`
8. Click **"Save"**

**Do this next:**
9. Go to **"Variables"** tab
10. Click **"+ New Variable"** and add:

    **Variable 1:**
    - Name: `HASURA_GRAPHQL_DATABASE_URL`
    - Value: Paste the DATABASE_URL from PostgreSQL (same as before)
    - Click "Add"

    **Variable 2:**
    - Name: `HASURA_GRAPHQL_ADMIN_SECRET`
    - Value: `myhasurasecret`
    - Click "Add"

    **Variable 3:**
    - Name: `HASURA_GRAPHQL_ENABLE_CONSOLE`
    - Value: `true`
    - Click "Add"

11. Railway will auto-redeploy

**Wait:** 2-3 minutes for Hasura to deploy

**Do this next:**
12. Go to **"Settings"** tab of Hasura service
13. Click **"Generate Domain"** button
14. **COPY THE URL** (e.g., `https://hasura-production-xxxx.up.railway.app`)

---

## 🔗 STEP 9: Connect Frontend to Backend

**Do this:**
1. Go back to **Vercel** (vercel.com)
2. Click on your project
3. Go to **"Settings"** (top menu)
4. Click **"Environment Variables"** (left sidebar)
5. Click **"+ Add"** button

**Add these variables one by one:**

**Variable 1:**
- Key: `NUXT_PUBLIC_API_URL`
- Value: Paste your **backend URL** from Railway (from Step 7)
- Environment: Select **Production, Preview, Development** (all three)
- Click "Save"

**Variable 2:**
- Key: `NUXT_PUBLIC_HASURA_URL`
- Value: Paste your **Hasura URL** from Railway (from Step 8) + `/v1/graphql`
  - Example: `https://hasura-production-xxxx.up.railway.app/v1/graphql`
- Environment: Select **Production, Preview, Development** (all three)
- Click "Save"

**Check:** You should see 2 environment variables listed

---

## 🔄 STEP 10: Redeploy Frontend

**Do this:**
1. In Vercel, go to **"Deployments"** tab (top menu)
2. Find your latest deployment
3. Click the **"..."** (three dots) menu
4. Click **"Redeploy"**
5. Click **"Redeploy"** again to confirm

**Wait:** 2-3 minutes for redeployment

**Check:** Deployment should show "Ready" with green checkmark

---

## ✅ STEP 11: Test Your App

**Do this:**
1. Click on your deployment in Vercel
2. Click the **"Visit"** button (or use your original link)
3. Your app should open in a new tab

**Check these:**
- ✅ App loads (you see the homepage)
- ✅ Can navigate pages
- ✅ Can register/login (test with a new account)
- ✅ Can view recipes

**If something doesn't work:**
- Check browser console (F12) for errors
- Check Vercel deployment logs
- Check Railway service logs

---

## 🎁 STEP 12: Get Your Final Link

**Do this:**
1. In Vercel, go to your project
2. Look at the top - you'll see your domain
3. It looks like: `https://food-recipes-app-xxxxx.vercel.app`
4. **This is your link!**

**Optional - Custom Domain:**
- Go to **Settings** → **Domains**
- Add your own domain (if you have one)

---

## 📧 STEP 13: Share with Your Boss

**Do this:**
1. Copy your Vercel link: `https://your-app.vercel.app`
2. Send it to your boss via:
   - Email
   - Slack/Teams
   - WhatsApp
   - Or just show them on your phone!

**Say:** "Here's the deployed app: [your link]"

---

## 🎉 DONE!

Your app is now live and accessible from anywhere in the world!

---

## 🆘 Troubleshooting

**Problem: Build fails in Vercel**
- **Fix:** Check Root Directory is `frontend/nuxt-app`
- **Fix:** Check build logs in Vercel for specific errors

**Problem: App loads but shows errors**
- **Fix:** Check environment variables are set correctly
- **Fix:** Check backend is running in Railway
- **Fix:** Check browser console (F12) for specific errors

**Problem: Can't login/register**
- **Fix:** Check backend URL in environment variables
- **Fix:** Check backend logs in Railway

**Problem: Database errors**
- **Fix:** Check DATABASE_URL is correct in Railway backend service
- **Fix:** Make sure PostgreSQL is running

---

## 📝 Quick Reference

**Your Links:**
- Frontend: `https://your-app.vercel.app` (from Vercel)
- Backend: `https://backend-xxxx.railway.app` (from Railway)
- Hasura: `https://hasura-xxxx.railway.app` (from Railway)

**Where to find things:**
- Vercel Dashboard: https://vercel.com/dashboard
- Railway Dashboard: https://railway.app/dashboard

---

**Follow these steps in order, and you'll have your app deployed! 🚀**

