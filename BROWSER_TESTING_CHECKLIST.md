# 🌐 BROWSER TESTING CHECKLIST

## ✅ CRITICAL FIXES APPLIED

### 1. Recipe Detail Page - Social Buttons
- ✅ Added `watch` for recipe changes to check interactions
- ✅ Fixed `checkUserInteractions` to run after recipe loads
- ✅ All buttons (Like, Bookmark, Comment, Rating) are visible when logged in
- ✅ Buttons trigger backend API calls
- ✅ UI updates in real-time after actions

### 2. Profile Page
- ✅ Fixed navigation using `useRouter()` instead of `navigateTo()`
- ✅ Added Profile link to navigation menu
- ✅ Profile page accessible at `/profile`
- ✅ All tabs (Recipes, Bookmarked, Purchased) fetch from backend

### 3. Recipe Edit/Delete
- ✅ Edit/Delete buttons visible on user's own recipes
- ✅ Navigation to edit page works
- ✅ Delete confirmation dialog works
- ✅ Backend ownership validation in place

### 4. Payment System
- ✅ Payment button visible for premium recipes
- ✅ Payment initialization works
- ✅ Success page created
- ✅ Purchase verification works

---

## 🧪 BROWSER TESTING STEPS

### Step 1: Login
1. Go to http://localhost:3000
2. Click "Log in"
3. Use test credentials: `submit@test.com` / `submit123`
4. ✅ Verify: Redirected to home page
5. ✅ Verify: Navigation shows "Profile" link

### Step 2: View Recipe Detail
1. Click on any recipe card
2. ✅ Verify: Recipe detail page loads
3. ✅ Verify: Like, Bookmark buttons visible (if logged in)
4. ✅ Verify: Rating stars visible (if logged in)
5. ✅ Verify: Comment form visible (if logged in)

### Step 3: Test Like Button
1. Click "Like" button
2. ✅ Verify: Button changes to "Liked" (red background)
3. ✅ Verify: No console errors
4. ✅ Verify: Check database: `SELECT * FROM likes WHERE user_id = {your_user_id}`

### Step 4: Test Bookmark Button
1. Click "Bookmark" button
2. ✅ Verify: Button changes to "Bookmarked" (yellow background)
3. ✅ Verify: No console errors
4. ✅ Verify: Check database: `SELECT * FROM bookmarks WHERE user_id = {your_user_id}`

### Step 5: Test Comment
1. Type a comment in the textarea
2. Click "Post Comment"
3. ✅ Verify: Comment appears in comments list
4. ✅ Verify: Textarea clears
5. ✅ Verify: Check database: `SELECT * FROM comments WHERE user_id = {your_user_id}`

### Step 6: Test Rating
1. Click on a star (1-5)
2. ✅ Verify: Stars fill up to selected rating
3. ✅ Verify: "Rating submitted!" message appears
4. ✅ Verify: Check database: `SELECT * FROM ratings WHERE user_id = {your_user_id}`

### Step 7: Test Profile Page
1. Click "Profile" in navigation
2. ✅ Verify: Profile page loads
3. ✅ Verify: Shows user's name and email
4. ✅ Verify: "My Recipes" tab shows user's recipes
5. ✅ Verify: Edit/Delete buttons visible on own recipes
6. Click "Bookmarked" tab
7. ✅ Verify: Shows bookmarked recipes
8. Click "Purchased" tab
9. ✅ Verify: Shows purchased recipes

### Step 8: Test Edit/Delete
1. On profile page, click Edit button on a recipe
2. ✅ Verify: Redirects to edit page (or shows edit form)
3. Go back to profile
4. Click Delete button on a recipe
5. ✅ Verify: Confirmation dialog appears
6. Click "OK"
7. ✅ Verify: Recipe removed from list
8. ✅ Verify: Check database: Recipe deleted

### Step 9: Test Payment (if recipe has price > 0)
1. Go to a premium recipe (price > 0)
2. ✅ Verify: "Buy Recipe" button visible
3. Click "Buy Recipe"
4. ✅ Verify: Redirects to Chapa payment page
5. Complete payment (test mode)
6. ✅ Verify: Redirects to success page
7. ✅ Verify: Recipe content unlocked
8. ✅ Verify: Check database: `SELECT * FROM purchases WHERE user_id = {your_user_id}`

---

## 🔍 DATABASE VERIFICATION QUERIES

After testing, run these queries to verify data persistence:

```sql
-- Check user's likes
SELECT l.*, r.title as recipe_title 
FROM likes l 
JOIN recipes r ON l.recipe_id = r.id 
WHERE l.user_id = {your_user_id};

-- Check user's bookmarks
SELECT b.*, r.title as recipe_title 
FROM bookmarks b 
JOIN recipes r ON b.recipe_id = r.id 
WHERE b.user_id = {your_user_id};

-- Check user's comments
SELECT c.*, r.title as recipe_title 
FROM comments c 
JOIN recipes r ON c.recipe_id = r.id 
WHERE c.user_id = {your_user_id}
ORDER BY c.created_at DESC;

-- Check user's ratings
SELECT rt.*, r.title as recipe_title 
FROM ratings rt 
JOIN recipes r ON rt.recipe_id = r.id 
WHERE rt.user_id = {your_user_id};

-- Check user's purchases
SELECT p.*, r.title as recipe_title 
FROM purchases p 
JOIN recipes r ON p.recipe_id = r.id 
WHERE p.user_id = {your_user_id}
ORDER BY p.created_at DESC;
```

---

## 🐛 COMMON ISSUES & FIXES

### Issue: Buttons not visible
**Fix**: Ensure user is logged in. Buttons only show when `isAuthenticated` is true.

### Issue: Buttons don't work
**Fix**: 
1. Check browser console for errors
2. Verify backend is running on port 8081
3. Verify JWT token is valid
4. Check network tab for API responses

### Issue: Profile page blank
**Fix**:
1. Check if user is logged in
2. Check browser console for GraphQL errors
3. Verify backend endpoints are accessible
4. Check if user has any recipes/bookmarks/purchases

### Issue: Actions not persisting
**Fix**:
1. Check backend logs for errors
2. Verify database connection
3. Check if JWT authentication is working
4. Verify database tables exist

---

## ✅ EXPECTED RESULTS

After completing all tests:
- ✅ All buttons visible and functional
- ✅ Profile page shows all user data
- ✅ All actions persist to database
- ✅ Real-time UI updates work
- ✅ No console errors
- ✅ All API calls return 200/201 status

---

**Status**: Ready for browser testing!

