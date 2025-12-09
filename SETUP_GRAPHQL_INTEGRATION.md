# GraphQL Integration Setup Guide

## ✅ Home Page Updated to Use GraphQL

The home page (`frontend/nuxt-app/pages/home.vue`) has been updated to use **GraphQL through Hasura** as per requirements:

### GraphQL Queries (Vue Apollo)
- ✅ Recipes query using `gql` and `useAsyncQuery`
- ✅ Categories query using `gql` and `useAsyncQuery`
- ✅ Connected to Hasura GraphQL endpoint: `http://localhost:8080/v1/graphql`

### Requirements Compliance
- ✅ **Vue Apollo** - Using `@vue/apollo-composable` and `useAsyncQuery`
- ✅ **Hasura** - GraphQL queries go through Hasura
- ✅ **GraphQL** - All data fetching uses GraphQL queries

---

## 🔧 Hasura Configuration Steps

### 1. Ensure Hasura is Running
```powershell
cd docker
docker-compose up -d
```

### 2. Track Tables in Hasura
```powershell
cd backend
.\configure_hasura.ps1
```

### 3. Configure Relationships
```powershell
.\configure_hasura_relationships.ps1
```

This creates:
- `recipes.user` - Relationship to users table
- `recipes.category` - Relationship to categories table
- `users.recipes` - Array relationship (one user has many recipes)

### 4. Configure Computed Fields
```powershell
.\configure_computed_fields.ps1
```

### 5. Configure Permissions
```powershell
.\configure_hasura_permissions.ps1
```

---

## 📋 GraphQL Schema (Hasura Auto-Generated)

### Recipes Query
```graphql
query GetRecipes {
  recipes(order_by: { created_at: desc }) {
    id
    title
    description
    thumbnail_url
    price
    preparation_time
    created_at
    category_id
    user_id
    user {
      id
      name
    }
    category {
      id
      name
    }
  }
}
```

### Categories Query
```graphql
query GetCategories {
  categories(order_by: { name: asc }) {
    id
    name
    image_url
  }
}
```

---

## ✅ Current Implementation

### Frontend (Home Page)
- ✅ Uses Vue Apollo (`useAsyncQuery`)
- ✅ GraphQL queries via Hasura
- ✅ Real-time data from database
- ✅ All filters working (client-side)

### Backend
- ✅ All tables tracked in Hasura
- ✅ Foreign keys create relationships automatically
- ✅ Computed fields configured
- ✅ Permissions script ready

---

## 🚀 Running the Full Stack

### Start All Services
1. **Docker Services** (PostgreSQL + Hasura):
   ```powershell
   cd docker
   docker-compose up -d
   ```

2. **Backend** (Go API):
   ```powershell
   cd backend
   go run main.go
   ```

3. **Frontend** (Nuxt.js):
   ```powershell
   cd frontend/nuxt-app
   npm run dev
   ```

### Configure Hasura (One-time setup)
```powershell
cd backend
.\configure_hasura.ps1
.\configure_hasura_relationships.ps1
.\configure_computed_fields.ps1
.\configure_hasura_permissions.ps1
```

---

## 🧪 Testing GraphQL

### Test in Hasura Console
1. Open: http://localhost:8080
2. Go to GraphQL tab
3. Run query:
```graphql
query {
  recipes {
    id
    title
    user {
      name
    }
  }
}
```

### Test in Browser
1. Open: http://localhost:3000/home
2. Should see recipes loaded via GraphQL
3. Check browser console for any errors

---

## 📝 Notes

- **GraphQL Endpoint**: `http://localhost:8080/v1/graphql`
- **Apollo Client**: Configured in `plugins/apollo.client.ts`
- **Authentication**: JWT token sent via Authorization header
- **Relationships**: Auto-detected from foreign keys, or configured via script

---

## ✅ Status

**Home page is now using GraphQL through Hasura as required!**

- ✅ Vue Apollo integration
- ✅ GraphQL queries
- ✅ Hasura as GraphQL engine
- ✅ All requirements met

