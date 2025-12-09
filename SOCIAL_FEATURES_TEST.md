# ✅ SOCIAL FEATURES INTEGRATION TEST

## 🧪 TEST EXECUTION

**Test Date**: December 9, 2025  
**Purpose**: Verify likes, comments, ratings, bookmarks persist in database

---

## 📊 TEST SETUP

### Test User Created:
- Email: social@test.com
- Password: social123
- JWT Token: ✅ Obtained

---

## 📊 TEST 1: LIKE FEATURE

### API Call:
```http
POST http://localhost:8081/recipes/1/like
Authorization: Bearer {JWT_TOKEN}
```

### Expected Database Action:
```sql
INSERT INTO likes (user_id, recipe_id, created_at)
VALUES (45, 1, NOW())
ON CONFLICT (user_id, recipe_id) DO NOTHING
```

### Result:
- Status: (Testing...)
- Database INSERT: (Verifying...)

---

## 📊 TEST 2: COMMENT FEATURE

### API Call:
```http
POST http://localhost:8081/recipes/1/comments
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json

{
  "content": "This recipe looks amazing!"
}
```

### Expected Database Action:
```sql
INSERT INTO comments (user_id, recipe_id, content, created_at)
VALUES (45, 1, 'This recipe looks amazing!', NOW())
RETURNING id
```

### Result:
- Status: (Testing...)
- Database INSERT: (Verifying...)

---

## 📊 TEST 3: RATING FEATURE

### API Call:
```http
POST http://localhost:8081/recipes/1/rate
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json

{
  "rating": 5
}
```

### Expected Database Action:
```sql
INSERT INTO ratings (user_id, recipe_id, rating, created_at)
VALUES (45, 1, 5, NOW())
ON CONFLICT (user_id, recipe_id)
DO UPDATE SET rating = 5, updated_at = NOW()
```

### Result:
- Status: (Testing...)
- Database INSERT: (Verifying...)

---

## 📊 TEST 4: BOOKMARK FEATURE

### API Call:
```http
POST http://localhost:8081/recipes/1/bookmark
Authorization: Bearer {JWT_TOKEN}
```

### Expected Database Action:
```sql
INSERT INTO bookmarks (user_id, recipe_id, created_at)
VALUES (45, 1, NOW())
ON CONFLICT (user_id, recipe_id) DO NOTHING
```

### Result:
- Status: (Testing...)
- Database INSERT: (Verifying...)

---

## 🔍 DATABASE VERIFICATION

### Query to Check Likes:
```graphql
{
  likes(where: {recipe_id: {_eq: 1}}) {
    id
    user_id
    recipe_id
    created_at
  }
}
```

### Query to Check Comments:
```graphql
{
  comments(where: {recipe_id: {_eq: 1}}) {
    id
    user_id
    content
    created_at
  }
}
```

### Query to Check Ratings:
```graphql
{
  ratings(where: {recipe_id: {_eq: 1}}) {
    id
    user_id
    rating
    created_at
  }
}
```

---

## ✅ EXPECTED RESULTS

After all tests:
- ✅ 1 like in `likes` table
- ✅ 1 comment in `comments` table
- ✅ 1 rating in `ratings` table
- ✅ 1 bookmark in `bookmarks` table

All for recipe_id = 1, user_id = 45

---

## 📝 VERIFICATION CHECKLIST

- [ ] Like button inserts into database
- [ ] Unlike button deletes from database
- [ ] Comment form saves to database
- [ ] Rating stars save to database
- [ ] Bookmark button inserts into database
- [ ] Profile shows user's likes
- [ ] Profile shows user's bookmarks
- [ ] Profile shows user's comments
- [ ] All data persists after page refresh

---

**Status**: Running tests...

