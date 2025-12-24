# 🍳 Food Recipes App - Complete Project Guide
## 5-Day Learning Plan from Scratch to Expert

---

## 📋 Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture & Technology Stack](#architecture--technology-stack)
3. [Folder Structure Explained](#folder-structure-explained)
4. [Day 1: Backend Foundation](#day-1-backend-foundation)
5. [Day 2: Database & Models](#day-2-database--models)
6. [Day 3: API Handlers & Business Logic](#day-3-api-handlers--business-logic)
7. [Day 4: Frontend Foundation](#day-4-frontend-foundation)
8. [Day 5: Integration & Advanced Features](#day-5-integration--advanced-features)
9. [Key Concepts Explained](#key-concepts-explained)

---

## 🎯 Project Overview

### What This Project Does
This is a **Food Recipes Application** where users can:
- Register and login
- Create, view, edit, and delete recipes
- Like, bookmark, and comment on recipes
- Purchase premium recipes
- Upload images for recipes
- Rate recipes
- View their profile with bookmarks and purchases

### Technology Stack
- **Backend**: Go (Golang) - REST API
- **Database**: PostgreSQL
- **GraphQL**: Hasura (for advanced queries)
- **Frontend**: Nuxt.js 3 (Vue.js framework)
- **Styling**: Tailwind CSS
- **Authentication**: JWT (JSON Web Tokens)
- **Payment**: Chapa (Ethiopian payment gateway)

---

## 🏗️ Architecture & Technology Stack

### System Architecture
```
┌─────────────┐
│   Browser   │ (User Interface)
└──────┬──────┘
       │ HTTP/GraphQL
       ▼
┌─────────────┐
│  Nuxt.js    │ (Frontend - Port 3000)
│  (Vue.js)   │
└──────┬──────┘
       │
       ├──────────────┐
       ▼              ▼
┌─────────────┐  ┌─────────────┐
│  Go Backend │  │   Hasura    │
│  (Port 8081)│  │ (Port 8080) │
└──────┬──────┘  └──────┬──────┘
       │                │
       └────────┬───────┘
                ▼
         ┌─────────────┐
         │ PostgreSQL  │
         │ (Port 5433) │
         └─────────────┘
```

### Why This Architecture?
1. **Go Backend**: Fast, efficient, handles concurrent requests well
2. **Hasura**: Auto-generates GraphQL API from database, saves development time
3. **PostgreSQL**: Reliable relational database for complex data relationships
4. **Nuxt.js**: Server-side rendering, better SEO, faster page loads

---

## 📁 Folder Structure Explained

```
food-recipes-app/
│
├── backend/                    # Go backend server
│   ├── main.go                 # Entry point, route definitions
│   ├── go.mod                  # Go dependencies
│   ├── handlers/               # Request handlers (business logic)
│   │   ├── auth.go            # Login, signup, JWT
│   │   ├── recipes.go         # Recipe CRUD operations
│   │   ├── recipes_handler.go # Additional recipe handlers
│   │   ├── profile.go         # User profile, bookmarks, purchases
│   │   ├── payment.go         # Payment processing
│   │   ├── social.go          # Likes, comments, ratings
│   │   ├── social_check.go    # Check if user liked/bookmarked
│   │   ├── extra.go           # File uploads, Hasura actions
│   │   └── events.go          # Hasura event triggers
│   ├── models/                 # Data structures
│   │   ├── user.go            # User model
│   │   ├── recipe.go          # Recipe, Category models
│   │   └── social.go          # Like, Comment, Rating models
│   ├── utils/                  # Helper functions
│   │   └── jwt.go             # JWT token generation/validation
│   ├── migrations/             # Database schema
│   │   ├── V1__create_core_tables.sql
│   │   ├── V2__create_recipe_images.sql
│   │   ├── V3__create_social_tables.sql
│   │   ├── V4__create_ratings_table.sql
│   │   └── V5__add_triggers.sql
│   ├── uploads/                # Uploaded images
│   └── *.ps1                   # PowerShell setup scripts
│
├── frontend/
│   └── nuxt-app/               # Nuxt.js application
│       ├── pages/              # Routes (auto-generated)
│       │   ├── index.vue       # Home page
│       │   ├── login.vue        # Login page
│       │   ├── register.vue     # Signup page
│       │   ├── create.vue       # Create recipe
│       │   ├── profile.vue      # User profile
│       │   └── recipes/
│       │       ├── [id].vue    # Recipe detail (dynamic route)
│       │       └── [id]/edit.vue # Edit recipe
│       ├── layouts/             # Page layouts
│       │   ├── default.vue     # Default layout (with nav)
│       │   └── blank.vue       # No navigation
│       ├── composables/         # Reusable logic
│       │   └── useAuth.js      # Authentication composable
│       ├── plugins/             # Nuxt plugins
│       │   ├── apollo.client.ts # GraphQL client setup
│       │   └── jwt-interceptor.client.js # JWT handling
│       ├── nuxt.config.ts      # Nuxt configuration
│       └── package.json        # Frontend dependencies
│
└── docker/
    └── docker-compose.yml      # Docker setup (Postgres + Hasura)
```

---

## 📅 Day 1: Backend Foundation

### Goal: Understand Go backend setup, main.go, and basic routing

---

### 1.1 Understanding `backend/main.go` - Line by Line

```go
package main  // Every Go file starts with package declaration
```

**Explanation**: `package main` tells Go this is an executable program (not a library).

```go
import (
	"foodrecipes/handlers"  // Our custom handlers package
	"log"                   // For logging
	"net/http"              // HTTP server and client
	"os"                    // Operating system interface
	"strings"               // String manipulation

	"github.com/jmoiron/sqlx"  // SQL database library
	_ "github.com/lib/pq"      // PostgreSQL driver (underscore = import for side effects only)
)
```

**Explanation**: 
- `import` brings in packages we need
- `_` before `github.com/lib/pq` means "import this but don't use it directly" - it registers the PostgreSQL driver

```go
func main() {
	// This is the entry point - runs when program starts
```

**Explanation**: `main()` is special - Go automatically calls this function when the program runs.

```go
	// Database connection string
	connStr := "postgres://fooduser:foodpass@127.0.0.1:5433/foodrecipes?sslmode=disable"
	if envConn := os.Getenv("DATABASE_URL"); envConn != "" {
		connStr = envConn
	}
```

**Explanation**:
- `connStr` is the database connection string
- Format: `postgres://username:password@host:port/database?options`
- `os.Getenv("DATABASE_URL")` checks for environment variable (for production)
- If environment variable exists, use it instead (more secure)

```go
	// Connect to Postgres
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")
```

**Explanation**:
- `sqlx.Connect()` tries to connect to database
- `err` will be `nil` if successful, or contain error if failed
- `log.Fatalf()` logs error and stops program (can't run without database)
- `log.Println()` logs success message

```go
	// Assign db to handlers.DB
	handlers.DB = db
```

**Explanation**: Makes database connection available to all handler functions (shared resource).

---

### 1.2 Understanding Route Registration

```go
	// Register login route
	http.HandleFunc("/login", handlers.LoginHandler)
```

**Explanation**:
- `http.HandleFunc()` registers a route
- `"/login"` is the URL path
- `handlers.LoginHandler` is the function that handles requests to `/login`
- When someone visits `http://localhost:8081/login`, `LoginHandler` runs

```go
	// Register file upload route
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.UploadFileHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
```

**Explanation**:
- This uses an **anonymous function** (function without a name)
- `r.Method` checks HTTP method (GET, POST, PUT, DELETE)
- Only allows POST requests
- `http.Error()` sends error response

---

### 1.3 Understanding Complex Routing

```go
	// Register recipe routes (Protected)
	http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AuthMiddleware(handlers.CreateRecipeHandler)(w, r)
		case http.MethodGet:
			handlers.GetRecipesHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
```

**Explanation**:
- `switch` statement checks HTTP method
- `POST /recipes` → Create recipe (requires authentication)
- `GET /recipes` → Get all recipes (public)
- `AuthMiddleware` wraps handler to check if user is logged in

**Middleware Pattern**:
```go
AuthMiddleware(handler)(w, r)
```
This is equivalent to:
```go
// AuthMiddleware returns a function
authenticatedHandler := AuthMiddleware(handler)
// Then we call that function
authenticatedHandler(w, r)
```

---

### 1.4 Understanding CORS Middleware

```go
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}
```

**Explanation**:
- **CORS** = Cross-Origin Resource Sharing
- Browser blocks requests from `localhost:3000` (frontend) to `localhost:8081` (backend) without CORS headers
- `Access-Control-Allow-Origin: *` allows any origin (for development)
- `OPTIONS` is a preflight request - browser asks "can I make this request?"
- `next.ServeHTTP(w, r)` passes request to next handler in chain

---

### 1.5 Starting the Server

```go
	log.Println("Starting Food Recipes Backend on :8081...")
	log.Fatal(http.ListenAndServe(":8081", corsMiddleware(http.DefaultServeMux)))
```

**Explanation**:
- `http.ListenAndServe(":8081", ...)` starts HTTP server on port 8081
- `corsMiddleware(http.DefaultServeMux)` wraps all routes with CORS
- `log.Fatal()` stops program if server fails to start
- Server runs forever until you stop it (Ctrl+C)

---

### Day 1 Summary
✅ **What you learned:**
- How Go programs start (`main()` function)
- How to connect to PostgreSQL database
- How to register HTTP routes
- How middleware works (CORS, Auth)
- How HTTP methods (GET, POST, PUT, DELETE) work

---

## 📅 Day 2: Database & Models

### Goal: Understand database schema, models, and data structures

---

### 2.1 Database Schema - `migrations/V1__create_core_tables.sql`

```sql
-- Create Users Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Explanation**:
- `CREATE TABLE IF NOT EXISTS` - creates table only if it doesn't exist
- `SERIAL` - auto-incrementing integer (1, 2, 3, ...)
- `PRIMARY KEY` - unique identifier for each row
- `TEXT` - string data type
- `NOT NULL` - field is required
- `UNIQUE` - no two users can have same email
- `DEFAULT CURRENT_TIMESTAMP` - automatically sets current time

**Why this structure?**
- `id` - unique identifier
- `name` - user's display name
- `email` - login identifier (unique)
- `password` - hashed password (never store plain text!)
- `avatar_url` - optional profile picture
- `created_at` - when account was created

```sql
-- Create Recipes Table
CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    description TEXT,
    preparation_time INT, -- in minutes
    price NUMERIC(10, 2) DEFAULT 0.00,
    thumbnail_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Explanation**:
- `REFERENCES users(id)` - **Foreign Key** - links recipe to user
- `ON DELETE CASCADE` - if user is deleted, delete their recipes too
- `ON DELETE SET NULL` - if category is deleted, set category_id to NULL (don't delete recipe)
- `NUMERIC(10, 2)` - decimal number with 10 total digits, 2 after decimal (e.g., 99.99)

**Database Relationships**:
```
users (1) ────< (many) recipes
categories (1) ────< (many) recipes
recipes (1) ────< (many) recipe_ingredients
recipes (1) ────< (many) recipe_steps
```

---

### 2.2 Go Models - `backend/models/user.go`

```go
package models

// User model for authentication and profile
type User struct {
	ID        int    `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	Name      string `db:"name" json:"name"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}
```

**Explanation**:
- `type User struct` - defines a structure (like a class in other languages)
- `db:"id"` - tells sqlx which database column to map to
- `json:"id"` - tells JSON encoder/decoder the field name
- `json:"-"` - means "don't include this field in JSON" (password should never be sent to frontend!)

**Why struct tags?**
- Database uses snake_case (`avatar_url`)
- Go uses PascalCase (`AvatarURL`)
- JSON uses camelCase (`avatarUrl`)
- Tags map between these different naming conventions

---

### 2.3 Recipe Model - `backend/models/recipe.go`

```go
type Recipe struct {
	ID              int           `db:"id" json:"id"`
	UserID          int           `db:"user_id" json:"user_id"`
	CategoryID      int           `db:"category_id" json:"category_id"`
	Title           string        `db:"title" json:"title"`
	Description     string        `db:"description" json:"description"`
	PreparationTime int           `db:"preparation_time" json:"preparation_time"`
	Price           float64       `db:"price" json:"price"`
	ThumbnailURL    string        `db:"thumbnail_url" json:"thumbnail_url"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	Images          []RecipeImage `json:"images"`  // Not in database - loaded separately
	FeaturedImageID int           `json:"featured_image_id"`
}
```

**Explanation**:
- `Images []RecipeImage` - slice (array) of images
- This field is NOT in database - it's populated by joining with `recipe_images` table
- `FeaturedImageID` - also computed, not stored directly

**Why separate images table?**
- One recipe can have many images
- Storing as array in one column is bad practice (hard to query, update)
- Separate table allows: easy queries, better performance, easier updates

---

### 2.4 Request/Response Models

```go
type CreateRecipeRequest struct {
	CategoryID      int                `json:"category_id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	PreparationTime int                `json:"preparation_time"`
	Price           float64            `json:"price"`
	ThumbnailURL    string             `json:"thumbnail_url"`
	Ingredients     []RecipeIngredient `json:"ingredients"`
	Steps           []RecipeStep       `json:"steps"`
	Images          []string           `json:"images"`
}
```

**Explanation**:
- This is what the frontend sends when creating a recipe
- Contains nested data (ingredients, steps, images)
- Backend will split this into multiple database inserts

---

### Day 2 Summary
✅ **What you learned:**
- Database table structure and relationships
- Foreign keys and cascading deletes
- Go structs and struct tags
- How to map database columns to Go structs
- Why we separate related data into different tables

---

## 📅 Day 3: API Handlers & Business Logic

### Goal: Understand how requests are processed, authentication, and business logic

---

### 3.1 Authentication Handler - `backend/handlers/auth.go`

#### Login Handler - Step by Step

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
```

**Explanation**:
- `w` = Response Writer (where we send response)
- `r` = Request (incoming HTTP request)
- Sets response header to JSON format

```go
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid request body"})
		return
	}
```

**Explanation**:
- `r.Body` - request body (JSON data from frontend)
- `json.NewDecoder().Decode()` - converts JSON to Go struct
- `&req` - pass pointer so Decode can modify the struct
- If error, send 400 Bad Request and return early

```go
	var user models.User
	err := DB.Get(&user, "SELECT id, name, email, password, COALESCE(avatar_url, '') as avatar_url FROM users WHERE email=$1", req.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid email or password"})
		return
	}
```

**Explanation**:
- `DB.Get()` - executes SQL query and stores result in `user` struct
- `$1` - parameter placeholder (prevents SQL injection)
- `COALESCE(avatar_url, '')` - if avatar_url is NULL, use empty string
- If user not found, return error (don't reveal if email exists - security)

```go
	// Check password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid email or password"})
		return
	}
```

**Explanation**:
- `bcrypt.CompareHashAndPassword()` - compares plain password with hashed password
- Passwords are NEVER stored in plain text
- If they don't match, return error

```go
	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Could not generate token"})
		return
	}

	user.Password = "" // Don't send password to frontend
	json.NewEncoder(w).Encode(LoginResponse{Token: token, User: user})
}
```

**Explanation**:
- Generate JWT token (like a temporary ID card)
- Clear password before sending response
- Send token and user info to frontend

---

### 3.2 JWT Token Generation - `backend/utils/jwt.go`

```go
var jwtSecret = []byte("your-secret-key")

func GenerateJWT(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(30 * time.Second).Unix(),
		"https://hasura.io/jwt/claims": jwt.MapClaims{
			"x-hasura-allowed-roles": []string{"user"},
			"x-hasura-default-role":  "user",
			"x-hasura-user-id":       strconv.Itoa(userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
```

**Explanation**:
- **JWT** = JSON Web Token (like a signed document)
- `claims` - data stored in token (user_id, email, expiration)
- `exp` - expiration time (30 seconds - should be longer in production!)
- `https://hasura.io/jwt/claims` - special claims for Hasura GraphQL
- `SignedString()` - signs token with secret (prevents tampering)

**How JWT Works**:
1. User logs in → Backend creates JWT
2. Frontend stores JWT (localStorage)
3. Frontend sends JWT with every request (in Authorization header)
4. Backend verifies JWT signature
5. If valid, extract user_id from token

---

### 3.3 Authentication Middleware

```go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing authorization header"})
			return
		}

		// Extract token (format: "Bearer <token>")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Verify token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
			return
		}

		// Extract user_id from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID := int(claims["user_id"].(float64))
		
		// Store user_id in request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)

		// Call next handler
		next(w, r)
	}
}
```

**Explanation**:
- **Middleware** - function that runs before the actual handler
- Checks if Authorization header exists
- Extracts and verifies JWT token
- If valid, stores user_id in request context
- If invalid, returns 401 Unauthorized

**Middleware Pattern**:
```
Request → AuthMiddleware → Handler
         (checks token)   (uses user_id)
```

---

### 3.4 Recipe Creation Handler

```go
func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user_id from context (set by AuthMiddleware)
	userID := r.Context().Value("user_id").(int)

	// Parse request body
	var req models.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	// Start database transaction
	tx, err := DB.Beginx()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if function exits early

	// Insert recipe
	var recipeID int
	err = tx.Get(&recipeID, `
		INSERT INTO recipes (user_id, category_id, title, description, preparation_time, price, thumbnail_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userID, req.CategoryID, req.Title, req.Description, req.PreparationTime, req.Price, req.ThumbnailURL)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert ingredients
	for _, ingredient := range req.Ingredients {
		tx.Exec(`
			INSERT INTO recipe_ingredients (recipe_id, name, quantity, unit)
			VALUES ($1, $2, $3, $4)
		`, recipeID, ingredient.Name, ingredient.Quantity, ingredient.Unit)
	}

	// Insert steps
	for _, step := range req.Steps {
		tx.Exec(`
			INSERT INTO recipe_steps (recipe_id, step_number, instruction, image_url)
			VALUES ($1, $2, $3, $4)
		`, recipeID, step.StepNumber, step.Instruction, step.ImageURL)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": recipeID})
}
```

**Explanation**:
- **Transaction** - ensures all database operations succeed or all fail
- If inserting recipe succeeds but ingredients fail, transaction rolls back
- `defer tx.Rollback()` - automatically rolls back if function exits early
- `tx.Commit()` - saves all changes

**Why Transactions?**
- Prevents partial data (recipe without ingredients)
- Ensures data consistency

---

### Day 3 Summary
✅ **What you learned:**
- How to parse JSON requests
- How to query database with sqlx
- How JWT authentication works
- How middleware protects routes
- How database transactions ensure data consistency
- How to handle errors properly

---

## 📅 Day 4: Frontend Foundation

### Goal: Understand Nuxt.js structure, Vue components, and GraphQL

---

### 4.1 Nuxt.js Structure

#### What is Nuxt.js?
- Framework built on Vue.js
- Provides: routing, server-side rendering, file-based routing
- **File-based routing**: Create file in `pages/` → automatic route

#### `pages/index.vue` - Home Page

```vue
<template>
  <div>
    <h1>Welcome to Food Recipes</h1>
  </div>
</template>

<script setup>
// This is the home page
// Accessible at: http://localhost:3000/
</script>
```

**Explanation**:
- `<template>` - HTML structure
- `<script setup>` - JavaScript logic
- File name `index.vue` = route `/`

#### `pages/login.vue` - Login Page

```vue
<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin">
      <input v-model="email" type="email" placeholder="Email" />
      <input v-model="password" type="password" placeholder="Password" />
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const email = ref('')
const password = ref('')

const handleLogin = async () => {
  const response = await fetch('http://localhost:8081/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: email.value, password: password.value })
  })
  
  const data = await response.json()
  
  if (data.token) {
    localStorage.setItem('token', data.token)
    navigateTo('/home')
  }
}
</script>
```

**Explanation**:
- `v-model="email"` - two-way data binding (input value ↔ variable)
- `@submit.prevent` - prevents form default submit, calls `handleLogin`
- `ref('')` - creates reactive variable (changes trigger UI updates)
- `fetch()` - makes HTTP request to backend
- `localStorage.setItem()` - stores token in browser
- `navigateTo()` - Nuxt function to change page

---

### 4.2 Composables - Reusable Logic

#### `composables/useAuth.js`

```javascript
export const useAuth = () => {
  const token = useState('token', () => null)
  
  const login = async (email, password) => {
    const response = await fetch('http://localhost:8081/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    })
    
    const data = await response.json()
    
    if (data.token) {
      token.value = data.token
      localStorage.setItem('token', data.token)
      return { success: true }
    }
    
    return { success: false, error: data.error }
  }
  
  const logout = () => {
    token.value = null
    localStorage.removeItem('token')
  }
  
  const isAuthenticated = computed(() => !!token.value)
  
  return {
    token,
    login,
    logout,
    isAuthenticated
  }
}
```

**Explanation**:
- **Composable** - reusable function (like a hook in React)
- `useState()` - Nuxt's reactive state (shared across components)
- `computed()` - derived value (automatically updates when token changes)
- `!!token.value` - converts to boolean (true if token exists)

**Why Composables?**
- Reuse authentication logic across multiple pages
- Centralized state management
- Easier to test and maintain

---

### 4.3 GraphQL with Apollo Client

#### `plugins/apollo.client.ts`

```typescript
export default defineNuxtPlugin((nuxtApp) => {
  const apolloClient = new ApolloClient({
    uri: 'http://localhost:8080/v1/graphql',
    cache: new InMemoryCache(),
    headers: {
      'x-hasura-admin-secret': 'myhasurasecret'
    }
  })
  
  nuxtApp.provide('apollo', apolloClient)
})
```

**Explanation**:
- **Apollo Client** - GraphQL client library
- Connects to Hasura GraphQL endpoint
- `InMemoryCache` - caches queries for performance
- Hasura auto-generates GraphQL API from database

**GraphQL vs REST**:
- REST: Multiple endpoints (`/recipes`, `/recipes/1`, `/recipes/1/comments`)
- GraphQL: Single endpoint, specify what data you need

**Example GraphQL Query**:
```graphql
query GetRecipes {
  recipes {
    id
    title
    description
    user {
      name
      email
    }
  }
}
```

Returns exactly the data you request, nothing more.

---

### 4.4 Dynamic Routes

#### `pages/recipes/[id].vue`

```vue
<template>
  <div v-if="recipe">
    <h1>{{ recipe.title }}</h1>
    <p>{{ recipe.description }}</p>
  </div>
</template>

<script setup>
const route = useRoute()
const recipeId = route.params.id  // Gets :id from URL

const { data: recipe } = await useFetch(`http://localhost:8081/recipes/${recipeId}`)
</script>
```

**Explanation**:
- `[id].vue` - dynamic route parameter
- `/recipes/123` → `route.params.id = "123"`
- `useFetch()` - Nuxt function for data fetching (like fetch but with SSR support)

---

### Day 4 Summary
✅ **What you learned:**
- Vue.js component structure (template, script, style)
- Nuxt.js file-based routing
- Reactive data with `ref()` and `computed()`
- How to make HTTP requests
- How to use composables for reusable logic
- GraphQL basics with Apollo Client

---

## 📅 Day 5: Integration & Advanced Features

### Goal: Understand payment integration, file uploads, and Hasura integration

---

### 5.1 File Upload Handler

#### `backend/handlers/extra.go` - UploadFileHandler

```go
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (file upload format)
	r.ParseMultipartForm(10 << 20) // 10MB limit
	
	// Get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	// Create unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	
	// Create destination file
	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	
	// Copy uploaded file to destination
	io.Copy(dst, file)
	
	// Return file URL
	json.NewEncoder(w).Encode(map[string]string{
		"url": fmt.Sprintf("http://localhost:8081/uploads/%s", filename),
	})
}
```

**Explanation**:
- `ParseMultipartForm()` - parses file upload request
- `10 << 20` - bit shift (10 * 2^20 = 10MB)
- `time.Now().UnixNano()` - unique timestamp for filename
- `io.Copy()` - copies file from upload to disk
- Returns URL so frontend can use the image

**Frontend Upload**:
```javascript
const formData = new FormData()
formData.append('file', fileInput.files[0])

const response = await fetch('http://localhost:8081/upload', {
  method: 'POST',
  body: formData
})

const { url } = await response.json()
// Use url in recipe creation
```

---

### 5.2 Payment Integration - Chapa

#### `backend/handlers/payment.go`

```go
func InitializePaymentHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	
	var req struct {
		RecipeID int     `json:"recipe_id"`
		Amount   float64 `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	// Create payment with Chapa API
	chapaReq := map[string]interface{}{
		"amount": req.Amount,
		"currency": "ETB",
		"email": "user@example.com",
		"first_name": "User",
		"tx_ref": fmt.Sprintf("recipe-%d-%d", req.RecipeID, userID),
		"callback_url": "http://localhost:8081/payment/callback",
	}
	
	// Make request to Chapa
	resp, _ := http.Post("https://api.chapa.co/v1/transaction/initialize", ...)
	
	// Return checkout URL to frontend
	json.NewEncoder(w).Encode(map[string]string{
		"checkout_url": "...",
	})
}
```

**Explanation**:
- **Payment Flow**:
  1. User clicks "Buy Recipe"
  2. Frontend calls `/payment/initialize`
  3. Backend creates payment with Chapa
  4. Backend returns checkout URL
  5. Frontend redirects user to Chapa checkout
  6. User pays on Chapa
  7. Chapa calls `/payment/callback` (webhook)
  8. Backend verifies payment and records purchase

**Webhook**:
```go
func PaymentCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Chapa sends payment status here
	// Verify payment is successful
	// Insert into purchases table
}
```

---

### 5.3 Hasura Actions

#### What are Hasura Actions?
- Custom business logic that Hasura can't do automatically
- Example: Authentication (password hashing, JWT generation)

#### `backend/handlers/auth.go` - HasuraLoginHandler

```go
func HasuraLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Hasura sends action payload
	var payload HasuraActionPayload
	json.NewDecoder(r.Body).Decode(&payload)
	
	// Extract login request
	req := payload.Input.Arg
	
	// Same login logic as regular handler
	// ... (check password, generate JWT)
	
	// Return Hasura-compatible response
	json.NewEncoder(w).Encode(HasuraLoginResponse{
		Token:  token,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	})
}
```

**Hasura Action Payload Format**:
```json
{
  "action": {
    "name": "login"
  },
  "input": {
    "arg": {
      "email": "user@example.com",
      "password": "password123"
    }
  }
}
```

**Why Hasura Actions?**
- Hasura can't hash passwords or generate JWTs
- Actions let you write custom Go code
- Hasura calls your Go endpoint
- Returns result to GraphQL client

---

### 5.4 Social Features - Likes, Comments, Ratings

#### Like Handler

```go
func ToggleLikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	recipeID := extractRecipeID(r.URL.Path)
	
	// Check if like exists
	var exists bool
	DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id=$1 AND recipe_id=$2)", 
		userID, recipeID)
	
	if exists {
		// Unlike - delete
		DB.Exec("DELETE FROM likes WHERE user_id=$1 AND recipe_id=$2", userID, recipeID)
	} else {
		// Like - insert
		DB.Exec("INSERT INTO likes (user_id, recipe_id) VALUES ($1, $2)", userID, recipeID)
	}
	
	json.NewEncoder(w).Encode(map[string]bool{"liked": !exists})
}
```

**Explanation**:
- **Toggle** - if liked, unlike; if not liked, like
- Simple insert/delete operation
- Returns new state (liked or not)

---

### 5.5 Event Triggers

#### What are Event Triggers?
- Automatically run code when database changes
- Example: When new recipe is created, send notification

#### `backend/handlers/events.go`

```go
func NewRecipeEventHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Event struct {
			Data struct {
				New struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
				} `json:"new"`
			} `json:"data"`
		} `json:"event"`
	}
	
	json.NewDecoder(r.Body).Decode(&payload)
	
	recipeID := payload.Event.Data.New.ID
	recipeTitle := payload.Event.Data.New.Title
	
	// Do something with new recipe
	// Example: Send email, update search index, etc.
	
	w.WriteHeader(http.StatusOK)
}
```

**How it works**:
1. User creates recipe → Insert into database
2. Hasura detects INSERT
3. Hasura calls `/events/new-recipe` endpoint
4. Your Go code runs (send notification, etc.)

---

### Day 5 Summary
✅ **What you learned:**
- How file uploads work (multipart form)
- How payment integration works (Chapa)
- How Hasura Actions extend GraphQL
- How social features work (likes, comments)
- How event triggers automate workflows

---

## 🔑 Key Concepts Explained

### 1. REST API vs GraphQL

**REST API** (Go Backend):
- Multiple endpoints: `/recipes`, `/recipes/1`, `/recipes/1/comments`
- Each endpoint returns fixed data structure
- Simple, easy to understand
- Used for: Authentication, file uploads, payments

**GraphQL** (Hasura):
- Single endpoint: `/v1/graphql`
- Client specifies what data it needs
- Flexible, efficient
- Used for: Complex queries, relationships

**Why Both?**
- REST for custom logic (auth, payments)
- GraphQL for data queries (recipes, users)

---

### 2. Authentication Flow

```
1. User enters email/password
   ↓
2. Frontend sends POST /login
   ↓
3. Backend checks password
   ↓
4. Backend generates JWT
   ↓
5. Frontend stores JWT (localStorage)
   ↓
6. Frontend sends JWT with every request (Authorization header)
   ↓
7. Backend middleware verifies JWT
   ↓
8. Backend extracts user_id from JWT
```

---

### 3. Database Relationships

**One-to-Many**:
- One user has many recipes
- One recipe has many ingredients
- One recipe has many steps

**Many-to-Many** (via junction table):
- Users can like many recipes
- Recipes can be liked by many users
- Junction table: `likes` (user_id, recipe_id)

---

### 4. Middleware Pattern

```
Request → CORS Middleware → Auth Middleware → Handler
         (adds headers)     (checks token)    (business logic)
```

Each middleware:
1. Receives request
2. Does something (add headers, check auth)
3. Passes to next middleware/handler

---

### 5. Error Handling

**Always check errors**:
```go
if err != nil {
    // Handle error
    return
}
```

**HTTP Status Codes**:
- `200` - Success
- `201` - Created
- `400` - Bad Request (invalid data)
- `401` - Unauthorized (not logged in)
- `403` - Forbidden (no permission)
- `404` - Not Found
- `500` - Internal Server Error

---

## 🎓 Final Exam Questions (Test Your Understanding)

1. **Why do we hash passwords?**
   - Answer: Security - if database is hacked, passwords are unreadable

2. **What is the purpose of JWT?**
   - Answer: Stateless authentication - server doesn't need to store sessions

3. **Why use transactions?**
   - Answer: Ensure data consistency - all operations succeed or all fail

4. **What is middleware?**
   - Answer: Code that runs before/after handlers (like filters)

5. **Why separate images into different table?**
   - Answer: One recipe can have many images - better database design

6. **What is the difference between REST and GraphQL?**
   - Answer: REST has multiple endpoints, GraphQL has one flexible endpoint

7. **How does file upload work?**
   - Answer: Multipart form data → server saves file → returns URL

8. **What are Hasura Actions?**
   - Answer: Custom business logic that Hasura calls when needed

---

## 📚 Additional Resources

### Go Learning
- [Go Tour](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)

### Vue.js/Nuxt.js
- [Vue.js Docs](https://vuejs.org/)
- [Nuxt.js Docs](https://nuxt.com/)

### Database
- [PostgreSQL Tutorial](https://www.postgresql.org/docs/)
- [SQL Tutorial](https://www.w3schools.com/sql/)

### GraphQL
- [GraphQL Docs](https://graphql.org/)
- [Hasura Docs](https://hasura.io/docs/)

---

## 🎯 Project Checklist

Before presenting to your boss, make sure you understand:

- [ ] How the project structure is organized
- [ ] How authentication works (JWT, middleware)
- [ ] How database relationships work
- [ ] How REST API endpoints are structured
- [ ] How GraphQL queries work
- [ ] How file uploads work
- [ ] How payment integration works
- [ ] How Hasura Actions and Events work
- [ ] How frontend components communicate with backend
- [ ] How error handling works

---

## 💡 Tips for Explaining to Your Boss

1. **Start with the big picture**: "This is a food recipes app with user authentication, recipe management, and payment features."

2. **Explain the architecture**: "We use Go for the backend API, PostgreSQL for the database, Hasura for GraphQL, and Nuxt.js for the frontend."

3. **Show the flow**: "When a user logs in, here's what happens step by step..."

4. **Explain the "why"**: "We use JWT because it's stateless and scalable. We use transactions to ensure data consistency."

5. **Be confident**: You built this! Even if you used help, you understand how it works now.

---

**Good luck with your presentation! 🚀**



