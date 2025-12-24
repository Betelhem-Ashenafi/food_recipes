# 🚀 Quick Reference Cheat Sheet
## For Your Boss Presentation

---

## 📊 Project Overview (30 seconds)

**What it is**: Food Recipes App - Users can create, share, and purchase recipes

**Tech Stack**:
- Backend: Go (Golang) - REST API on port 8081
- Database: PostgreSQL on port 5433
- GraphQL: Hasura on port 8080
- Frontend: Nuxt.js (Vue.js) on port 3000

**Key Features**:
- User authentication (JWT)
- Recipe CRUD operations
- Social features (likes, comments, ratings)
- Payment integration (Chapa)
- Image uploads

---

## 🗂️ Folder Structure (1 minute)

```
backend/
├── main.go              → Server entry point, routes
├── handlers/           → Business logic (auth, recipes, payments)
├── models/             → Data structures (User, Recipe)
├── utils/              → Helper functions (JWT)
└── migrations/         → Database schema

frontend/nuxt-app/
├── pages/              → Routes (auto-generated from files)
├── layouts/            → Page layouts
├── composables/        → Reusable logic (useAuth)
└── plugins/            → Apollo Client, JWT interceptor
```

---

## 🔐 Authentication Flow (2 minutes)

**Login Process**:
1. User sends email/password → `POST /login`
2. Backend checks password (bcrypt)
3. Backend generates JWT token
4. Frontend stores token in localStorage
5. Frontend sends token in `Authorization: Bearer <token>` header
6. Middleware verifies token on protected routes

**Key Code Locations**:
- Login: `backend/handlers/auth.go` → `LoginHandler()`
- JWT: `backend/utils/jwt.go` → `GenerateJWT()`
- Middleware: `backend/handlers/auth.go` → `AuthMiddleware()`

---

## 🗄️ Database Structure (2 minutes)

**Main Tables**:
- `users` - User accounts
- `recipes` - Recipe data
- `recipe_ingredients` - Recipe ingredients
- `recipe_steps` - Cooking steps
- `likes` - User likes (many-to-many)
- `bookmarks` - User bookmarks
- `comments` - Recipe comments
- `ratings` - Recipe ratings
- `purchases` - Paid recipes

**Relationships**:
- User → Recipes (1-to-many)
- Recipe → Ingredients (1-to-many)
- Recipe → Steps (1-to-many)
- User ↔ Recipes (many-to-many via likes/bookmarks)

**Key Code**: `backend/migrations/V1__create_core_tables.sql`

---

## 🔌 API Endpoints (2 minutes)

**Public Endpoints**:
- `POST /login` - User login
- `POST /signup` - User registration
- `GET /recipes` - List all recipes
- `GET /recipes/{id}` - Get recipe details
- `GET /categories` - List categories

**Protected Endpoints** (require JWT):
- `POST /recipes` - Create recipe
- `PUT /recipes/{id}` - Update recipe
- `DELETE /recipes/{id}` - Delete recipe
- `POST /recipes/{id}/like` - Like/unlike recipe
- `POST /recipes/{id}/bookmark` - Bookmark recipe
- `POST /payment/initialize` - Start payment

**Key Code**: `backend/main.go` → route definitions

---

## 💳 Payment Flow (1 minute)

1. User clicks "Buy Recipe"
2. Frontend calls `POST /payment/initialize`
3. Backend creates payment with Chapa API
4. Backend returns checkout URL
5. User redirected to Chapa payment page
6. User completes payment
7. Chapa calls `POST /payment/callback` (webhook)
8. Backend verifies payment
9. Backend records purchase in database

**Key Code**: `backend/handlers/payment.go`

---

## 📤 File Upload Flow (1 minute)

1. User selects image file
2. Frontend sends `POST /upload` with FormData
3. Backend saves file to `uploads/` directory
4. Backend generates unique filename (timestamp)
5. Backend returns file URL
6. Frontend uses URL in recipe creation

**Key Code**: `backend/handlers/extra.go` → `UploadFileHandler()`

---

## 🔄 Hasura Integration (2 minutes)

**What is Hasura?**
- Auto-generates GraphQL API from database
- Provides real-time subscriptions
- Handles complex queries automatically

**Hasura Actions**:
- Custom business logic Hasura can't do
- Example: Authentication (password hashing, JWT)
- Hasura calls our Go endpoints: `/hasura/login`, `/hasura/signup`

**Event Triggers**:
- Automatically run code when database changes
- Example: New recipe created → send notification
- Hasura calls: `/events/new-recipe`

**Key Code**: 
- Actions: `backend/handlers/auth.go` → `HasuraLoginHandler()`
- Events: `backend/handlers/events.go` → `NewRecipeEventHandler()`

---

## 🎨 Frontend Structure (2 minutes)

**Pages** (auto-routing):
- `pages/index.vue` → `/` (home)
- `pages/login.vue` → `/login`
- `pages/register.vue` → `/register`
- `pages/create.vue` → `/create` (create recipe)
- `pages/recipes/[id].vue` → `/recipes/123` (dynamic route)

**Composables**:
- `useAuth()` - Authentication logic (login, logout, token)

**Plugins**:
- `apollo.client.ts` - GraphQL client setup
- `jwt-interceptor.client.js` - Adds JWT to requests

**Key Code**: `frontend/nuxt-app/pages/` and `composables/useAuth.js`

---

## 🔧 Key Technologies Explained

**Go (Golang)**:
- Fast, compiled language
- Great for concurrent requests
- Simple syntax
- Used for: Backend API server

**PostgreSQL**:
- Relational database
- ACID compliant (data integrity)
- Supports complex queries
- Used for: Data storage

**JWT (JSON Web Tokens)**:
- Stateless authentication
- Token contains user info
- Signed with secret key
- Used for: User authentication

**GraphQL**:
- Query language for APIs
- Client specifies what data it needs
- Single endpoint
- Used for: Complex data queries

**Vue.js/Nuxt.js**:
- Reactive frontend framework
- Component-based
- Server-side rendering
- Used for: User interface

---

## 🐛 Common Questions & Answers

**Q: Why Go instead of Node.js/Python?**
A: Go is faster, handles concurrency better, and compiles to a single binary.

**Q: Why PostgreSQL instead of MongoDB?**
A: We need relational data (users, recipes, ingredients). PostgreSQL handles relationships better.

**Q: Why both REST and GraphQL?**
A: REST for custom logic (auth, payments). GraphQL for flexible data queries.

**Q: How is security handled?**
A: Passwords hashed with bcrypt, JWT tokens for auth, SQL injection prevented with parameterized queries.

**Q: How does the frontend communicate with backend?**
A: HTTP requests (fetch) for REST API, Apollo Client for GraphQL queries.

**Q: What happens if database is down?**
A: Backend can't start (fails on connection). Frontend shows errors when API calls fail.

**Q: How are images stored?**
A: Files saved to `backend/uploads/` directory, served as static files via `/uploads/` route.

**Q: How does payment verification work?**
A: Chapa sends webhook to `/payment/callback` with payment status. Backend verifies and records purchase.

---

## 📝 Code Flow Examples

### Creating a Recipe:
```
1. User fills form on /create page
2. Frontend: POST /recipes with recipe data
3. Backend: AuthMiddleware checks JWT
4. Backend: CreateRecipeHandler() runs
5. Backend: Start database transaction
6. Backend: Insert recipe, ingredients, steps
7. Backend: Commit transaction
8. Backend: Return recipe ID
9. Frontend: Redirect to /recipes/{id}
```

### Liking a Recipe:
```
1. User clicks like button
2. Frontend: POST /recipes/{id}/like
3. Backend: AuthMiddleware checks JWT
4. Backend: ToggleLikeHandler() runs
5. Backend: Check if like exists
6. Backend: Insert or delete like
7. Backend: Return new like status
8. Frontend: Update UI (heart icon)
```

---

## 🎯 Presentation Tips

1. **Start with demo**: Show the app working
2. **Explain architecture**: Draw the diagram (Browser → Frontend → Backend → Database)
3. **Walk through one feature**: "Let me show you how login works..."
4. **Show code**: Open key files, explain what they do
5. **Answer questions**: Use this cheat sheet for quick answers

---

## 🔍 Key Files to Show

**Backend**:
- `backend/main.go` - Show routing
- `backend/handlers/auth.go` - Show login logic
- `backend/models/recipe.go` - Show data structures
- `backend/migrations/V1__create_core_tables.sql` - Show database schema

**Frontend**:
- `frontend/nuxt-app/pages/login.vue` - Show login page
- `frontend/nuxt-app/composables/useAuth.js` - Show auth logic
- `frontend/nuxt-app/pages/recipes/[id].vue` - Show recipe detail page

---

## ✅ Confidence Boosters

- You understand the architecture
- You can explain each component
- You know why each technology was chosen
- You can trace a request from frontend to database
- You understand security measures
- You can explain error handling

**You've got this! 💪**



