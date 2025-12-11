# Access Control Model Documentation

## Overview

This document describes the access control model for the Food Recipes application, detailing which features are publicly accessible and which require authentication.

---

## 🔓 Public Access (No Login Required)

The following features are available to all visitors without authentication:

### 1. Browse Recipes
- **Endpoint**: `GET /recipes`
- **Handler**: `GetRecipesHandler` (no AuthMiddleware)
- **Features**:
  - View all public recipes
  - Filter by category (`?category={id}`)
  - Filter by creator (`?creator={name}`)
  - Filter by title (`?title={search}`)
  - Filter by ingredient (`?ingredient={name}`)
  - Filter by preparation time (`?time={minutes}`)
- **Location**: `backend/main.go:84-86`

### 2. View Recipe Details
- **Endpoint**: `GET /recipes/{id}`
- **Handler**: `GetRecipeByIDHandler` (no AuthMiddleware)
- **Returns**: Full recipe details including title, description, images, etc.
- **Location**: `backend/main.go:198`

### 3. View Recipe Components
- **Ingredients**: `GET /recipes/{id}/ingredients` (public)
- **Steps**: `GET /recipes/{id}/steps` (public)
- **Images**: `GET /recipes/{id}/images` (public)
- **Location**: `backend/main.go:99-111, 173-177`

### 4. View Comments and Ratings
- **Comments**: `GET /recipes/{id}/comments` (public - view only)
- **Ratings**: `GET /recipes/{id}/rate` (public - view only)
- **Location**: `backend/main.go:128-130, 139-141`

### 5. Browse Categories
- **Endpoint**: `GET /categories`
- **Handler**: `GetCategoriesHandler` (no AuthMiddleware)
- **Location**: `backend/main.go:37`

---

## 🔐 Login/Signup Required

The following features require authentication via JWT token:

### 1. Purchase Paid Recipes
- **Initialize Payment**: `POST /payment/initialize`
  - **Handler**: `InitializePaymentHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:211`
- **Verify Payment**: `GET /payment/verify`
  - **Handler**: `VerifyPaymentHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:219`
- **Check Purchase**: `GET /recipes/{id}/purchase/check`
  - **Handler**: `CheckPurchaseHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:167`

### 2. Interact with Recipes
- **Like Recipe**: `POST /recipes/{id}/like`
  - **Handler**: `ToggleLikeHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:116`
- **Bookmark Recipe**: `POST /recipes/{id}/bookmark`
  - **Handler**: `ToggleBookmarkHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:122`
- **Comment on Recipe**: `POST /recipes/{id}/comments`
  - **Handler**: `PostCommentHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:132`
- **Rate Recipe**: `POST /recipes/{id}/rate`
  - **Handler**: `RateRecipeHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:143`
- **Check Like Status**: `GET /recipes/{id}/like/check`
  - **Handler**: `CheckLikeHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:151`
- **Check Bookmark Status**: `GET /recipes/{id}/bookmark/check`
  - **Handler**: `CheckBookmarkHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:159`

### 3. Manage Recipes
- **Create Recipe**: `POST /recipes`
  - **Handler**: `CreateRecipeHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:83`
- **Edit Recipe**: `PUT /recipes/{id}`
  - **Handler**: `EditRecipeHandler` (AuthMiddleware)
  - **Additional Check**: Verifies user owns the recipe
  - **Location**: `backend/main.go:200`
- **Delete Recipe**: `DELETE /recipes/{id}`
  - **Handler**: `DeleteRecipeHandler` (AuthMiddleware)
  - **Additional Check**: Verifies user owns the recipe
  - **Location**: `backend/main.go:202`
- **Upload Recipe Images**: `POST /recipes/{id}/images`
  - **Handler**: `UploadRecipeImagesHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:180`
- **Set Featured Image**: `POST /recipes/{id}/images/{imgId}/feature`
  - **Handler**: `FeatureRecipeImageHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:188`

### 4. Profile Management
- **View Bookmarks**: `GET /users/{id}/bookmarks`
  - **Handler**: `GetUserBookmarksHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:64`
- **View Purchases**: `GET /users/{id}/purchases`
  - **Handler**: `GetUserPurchasesHandler` (AuthMiddleware)
  - **Location**: `backend/main.go:71`

---

## Authentication Flow

### JWT Token Generation
1. User logs in via `POST /login` or `POST /hasura/login`
2. Backend validates credentials and generates JWT token
3. Token includes:
   - `user_id`: User identifier
   - `email`: User email
   - Hasura claims for GraphQL access

### Token Usage
- **REST API**: Send token in `Authorization: Bearer <token>` header
- **GraphQL (Hasura)**: Token automatically included via Apollo client
- **Validation**: `AuthMiddleware` validates token on protected routes

### Token Location
- **Backend Generation**: `backend/utils/jwt.go`
- **Backend Validation**: `backend/handlers/recipes.go` (AuthMiddleware)
- **Frontend Storage**: Cookie (`auth_token`)
- **Frontend Usage**: `frontend/nuxt-app/plugins/apollo.client.ts`

---

## Access Control Implementation

### Public Routes (No Middleware)
```go
// Examples from main.go
http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        handlers.GetRecipesHandler(w, r)  // No AuthMiddleware
    }
})
```

### Protected Routes (With Middleware)
```go
// Examples from main.go
http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        handlers.AuthMiddleware(handlers.CreateRecipeHandler)(w, r)  // AuthMiddleware
    }
})
```

### Ownership Verification
Some protected routes also verify ownership:
```go
// Example from EditRecipeHandler
// 1. Extract user ID from token (via AuthMiddleware)
// 2. Check recipe ownership
var ownerID int
err = DB.Get(&ownerID, "SELECT user_id FROM recipes WHERE id=$1", recipeID)
if ownerID != userID {
    http.Error(w, "Forbidden: not recipe owner", http.StatusForbidden)
    return
}
```

---

## Summary Table

| Feature | Endpoint | Auth Required | Ownership Check |
|---------|----------|---------------|-----------------|
| Browse Recipes | `GET /recipes` | ❌ No | N/A |
| View Recipe | `GET /recipes/{id}` | ❌ No | N/A |
| View Ingredients | `GET /recipes/{id}/ingredients` | ❌ No | N/A |
| View Steps | `GET /recipes/{id}/steps` | ❌ No | N/A |
| View Images | `GET /recipes/{id}/images` | ❌ No | N/A |
| View Comments | `GET /recipes/{id}/comments` | ❌ No | N/A |
| View Ratings | `GET /recipes/{id}/rate` | ❌ No | N/A |
| Browse Categories | `GET /categories` | ❌ No | N/A |
| Create Recipe | `POST /recipes` | ✅ Yes | N/A |
| Edit Recipe | `PUT /recipes/{id}` | ✅ Yes | ✅ Yes |
| Delete Recipe | `DELETE /recipes/{id}` | ✅ Yes | ✅ Yes |
| Like Recipe | `POST /recipes/{id}/like` | ✅ Yes | N/A |
| Bookmark Recipe | `POST /recipes/{id}/bookmark` | ✅ Yes | N/A |
| Comment on Recipe | `POST /recipes/{id}/comments` | ✅ Yes | N/A |
| Rate Recipe | `POST /recipes/{id}/rate` | ✅ Yes | N/A |
| Purchase Recipe | `POST /payment/initialize` | ✅ Yes | N/A |
| Verify Payment | `GET /payment/verify` | ✅ Yes | N/A |
| View Bookmarks | `GET /users/{id}/bookmarks` | ✅ Yes | N/A |
| View Purchases | `GET /users/{id}/purchases` | ✅ Yes | N/A |

---

## Frontend Implementation

### Public Access
- Home page (`/home`) - Browse recipes without login
- Recipe detail page (`/recipes/[id]`) - View full recipe details
- Category filtering - Available to all users

### Protected Access
- Create recipe page (`/create`) - Redirects to login if not authenticated
- Edit recipe page (`/recipes/[id]/edit`) - Requires authentication + ownership
- Profile pages - Require authentication
- Payment flow - Requires authentication

### Authentication Check
```typescript
// Example from frontend
const token = useCookie('auth_token');
onMounted(() => {
  if (!token.value) {
    router.push('/login');
  }
});
```

---

## Benefits of This Model

1. **Low Barrier to Entry**: Users can explore content without registration
2. **Encourages Registration**: Premium features require login
3. **User Engagement**: Social features (like, comment, rate) drive signups
4. **Content Protection**: Recipe creation/editing requires authentication
5. **Monetization**: Paid recipes require purchase (authentication required)

---

## Security Considerations

1. **JWT Validation**: All protected routes validate JWT tokens
2. **Ownership Checks**: Edit/delete operations verify recipe ownership
3. **Password Hashing**: User passwords are hashed with bcrypt
4. **CORS Configuration**: CORS middleware allows cross-origin requests (configure for production)

---

## Testing Access Control

### Test Public Access
```bash
# Should work without authentication
curl http://localhost:8081/recipes
curl http://localhost:8081/recipes/1
curl http://localhost:8081/categories
```

### Test Protected Access
```bash
# Should fail without token
curl -X POST http://localhost:8081/recipes

# Should work with token
curl -X POST http://localhost:8081/recipes \
  -H "Authorization: Bearer <token>"
```

---

## Conclusion

The access control model successfully balances:
- **Open access** for content discovery
- **Protected access** for user interactions and premium features
- **Ownership verification** for content management

This encourages user registration while maintaining an open browsing experience.




