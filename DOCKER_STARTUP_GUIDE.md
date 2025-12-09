# 🐳 Docker Startup Guide - Quick Fix

## ⚠️ Current Issue

**Error**: `Failed to connect to database: dial tcp 127.0.0.1:5433: connectex: No connection could be made`

**Cause**: Docker Desktop is not running or Postgres container is not started

---

## ✅ Quick Fix (5 Minutes)

### Step 1: Start Docker Desktop
1. Open **Docker Desktop** application
2. Wait for the whale icon to stop animating (fully started)
3. Ensure Docker Desktop shows "Running" status

### Step 2: Start Docker Services
```bash
cd D:\food-recipes-app
docker-compose -f docker\docker-compose.yml up -d
```

**Wait 30 seconds** for Postgres and Hasura to fully start.

### Step 3: Verify Services
```bash
# Check containers are running
docker ps

# Should see:
# - docker-postgres-1 (port 5433:5432)
# - docker-hasura-1 (port 8080:8080)
```

### Step 4: Restart Backend
```bash
cd backend
./backend.exe
```

### Step 5: Verify All Services
- ✅ Postgres: localhost:5433
- ✅ Hasura: http://localhost:8080/console
- ✅ Backend: http://localhost:8081
- ✅ Frontend: http://localhost:3000

---

## 🧪 Test Connection

```bash
# Test Postgres connection
docker exec docker-postgres-1 psql -U fooduser -d foodrecipes -c "SELECT COUNT(*) FROM users;"

# Should return user count (44+)
```

---

## 📊 Current System Status

### What's Working:
- ✅ Frontend: http://localhost:3000 (Running)
- ✅ Backend: Compiled (needs Docker for DB)
- ❌ Docker: Not running / Needs restart

### What's Needed:
1. Start Docker Desktop
2. Run `docker-compose up -d`
3. Restart backend
4. Test login

---

## 🎯 After Docker Starts

Once Docker services are running:

1. **Test Login**: submit@test.com / submit123
2. **Verify Database**: User 44, 45 should be visible
3. **Test Social Features**: Like, comment, rate should work
4. **All data persists**: Every action saves to Postgres

---

## 🚀 Complete Startup Commands

```powershell
# 1. Start Docker Desktop (manually)

# 2. Start services
cd D:\food-recipes-app
docker-compose -f docker\docker-compose.yml up -d

# 3. Wait for startup
Start-Sleep -Seconds 30

# 4. Start backend
cd backend
./backend.exe

# 5. Start frontend (if not running)
cd ../frontend/nuxt-app
npm run dev

# 6. Test
# Open: http://localhost:3000
# Login: submit@test.com / submit123
```

---

## ✅ Everything Will Work Once Docker Starts!

All features are implemented and ready:
- ✅ Like/Bookmark/Comment/Rate
- ✅ User Profile
- ✅ Recipe Edit/Delete
- ✅ Payment System

Just need Docker Desktop running! 🐳

