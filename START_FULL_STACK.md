# Full Stack Startup Guide

## ✅ Services Running

### 1. Docker Services (PostgreSQL + Hasura)
- **PostgreSQL**: `localhost:5433`
- **Hasura Console**: `http://localhost:8080` (Secret: `myhasurasecret`)
- **Status**: ✅ Running

### 2. Backend Server (Go)
- **URL**: `http://localhost:8081`
- **Status**: ✅ Running
- **Endpoints Available**:
  - `GET /recipes` - Browse recipes
  - `GET /categories` - Get categories
  - `POST /login` - Login
  - `POST /signup` - Signup
  - All other endpoints as documented

### 3. Frontend Server (Nuxt.js)
- **URL**: `http://localhost:3000`
- **Status**: ✅ Running
- **Pages**:
  - `/` - Welcome page
  - `/login` - Login page
  - `/register` - Registration page
  - `/home` - Home page (integrated with backend)

---

## 🚀 Access Your Application

### Frontend
- **Home Page**: http://localhost:3000
- **Login**: http://localhost:3000/login
- **Register**: http://localhost:3000/register
- **Home Feed**: http://localhost:3000/home

### Backend API
- **Base URL**: http://localhost:8081
- **Categories**: http://localhost:8081/categories
- **Recipes**: http://localhost:8081/recipes

### Hasura Console
- **URL**: http://localhost:8080
- **Admin Secret**: `myhasurasecret`

---

## 📋 Home Page Integration

The home page (`/home`) is now fully integrated with the backend:

### ✅ Features Working:
1. **Browse Recipes** - Fetches from `GET /recipes`
2. **Categories** - Fetches from `GET /categories`
3. **Search by Title** - Client-side filtering
4. **Filter by Category** - Client-side filtering
5. **Filter by Preparation Time** - Client-side filtering
6. **Sort Options** - Newest, Oldest, Title A-Z

### 🔄 Data Flow:
- **Recipes**: REST API → `http://localhost:8081/recipes`
- **Categories**: REST API → `http://localhost:8081/categories`
- **Real-time Updates**: Refresh page to see new recipes

---

## 🛠️ Manual Startup (If Needed)

### Start Docker Services
```powershell
cd docker
docker-compose up -d
```

### Start Backend
```powershell
cd backend
go run main.go
```

### Start Frontend
```powershell
cd frontend/nuxt-app
npm run dev
```

---

## 🧪 Test the Integration

1. **Open Browser**: http://localhost:3000
2. **Navigate to Home**: Click through to `/home` or login
3. **View Recipes**: Should see recipes from backend
4. **View Categories**: Should see categories displayed
5. **Test Filters**: Try searching, filtering by category, time
6. **Test Search**: Type in search bar to filter recipes

---

## 📝 Notes

- Backend runs on port **8081**
- Frontend runs on port **3000**
- PostgreSQL runs on port **5433**
- Hasura runs on port **8080**
- All services are running in separate windows/processes
- Home page uses REST API (not GraphQL) for better compatibility

---

## ✅ Status: FULL STACK RUNNING

All services are up and the home page is integrated with the backend!

**Next Steps:**
1. Open http://localhost:3000 in your browser
2. Navigate to `/home` to see recipes
3. Test all the filters and search functionality
4. Create an account to start adding recipes

