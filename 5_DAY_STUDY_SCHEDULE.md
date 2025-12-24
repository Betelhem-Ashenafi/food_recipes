# 📅 5-Day Study Schedule
## Prepare for Your Boss Presentation

---

## 🎯 How to Use This Schedule

- **Day 1-3**: Focus on backend (Go, database, API)
- **Day 4**: Focus on frontend (Vue.js, Nuxt.js)
- **Day 5**: Integration, testing, and practice presentation

**Time per day**: 2-3 hours
**Total time**: 10-15 hours

---

## 📆 Day 1: Backend Foundation (3 hours)

### Morning (1.5 hours)

#### 1. Read `COMPLETE_PROJECT_GUIDE.md` - Day 1 Section (30 min)
- Understand Go basics
- Understand `main.go` structure
- Understand routing

#### 2. Study `backend/main.go` Line by Line (45 min)
- Open the file
- Read each section
- Understand what each route does
- Write notes on:
  - What each route is for
  - Which routes need authentication
  - How middleware works

#### 3. Practice (15 min)
- Draw a diagram of all routes
- List which routes are public vs protected

### Afternoon (1.5 hours)

#### 4. Study CORS Middleware (30 min)
- Read the `corsMiddleware` function
- Understand why CORS is needed
- Understand preflight requests

#### 5. Study Database Connection (30 min)
- Understand connection string format
- Understand environment variables
- Understand error handling

#### 6. Practice (30 min)
- Try to explain to yourself:
  - How the server starts
  - How routes are registered
  - How middleware works
  - What happens when a request comes in

### ✅ Day 1 Checklist
- [ ] I understand how `main.go` works
- [ ] I can explain what each route does
- [ ] I understand middleware concept
- [ ] I can explain CORS
- [ ] I understand database connection

---

## 📆 Day 2: Database & Models (3 hours)

### Morning (1.5 hours)

#### 1. Read `COMPLETE_PROJECT_GUIDE.md` - Day 2 Section (30 min)
- Understand database schema
- Understand relationships
- Understand Go models

#### 2. Study Database Migrations (45 min)
- Read `backend/migrations/V1__create_core_tables.sql`
- Understand each table structure
- Understand foreign keys
- Understand `ON DELETE CASCADE` vs `ON DELETE SET NULL`

#### 3. Study Other Migrations (15 min)
- Read V2, V3, V4, V5 migrations
- Understand what each adds

### Afternoon (1.5 hours)

#### 4. Study Go Models (45 min)
- Read `backend/models/user.go`
- Read `backend/models/recipe.go`
- Read `backend/models/social.go`
- Understand struct tags (`db:`, `json:`)
- Understand why `json:"-"` for password

#### 5. Practice (45 min)
- Draw database relationship diagram
- List all tables and their relationships
- Explain why each field exists
- Practice explaining models to yourself

### ✅ Day 2 Checklist
- [ ] I understand database schema
- [ ] I can explain relationships (1-to-many, many-to-many)
- [ ] I understand Go structs and tags
- [ ] I can explain why password has `json:"-"`
- [ ] I understand foreign keys

---

## 📆 Day 3: API Handlers & Business Logic (3 hours)

### Morning (1.5 hours)

#### 1. Read `COMPLETE_PROJECT_GUIDE.md` - Day 3 Section (30 min)
- Understand authentication flow
- Understand JWT
- Understand middleware

#### 2. Study Authentication (45 min)
- Read `backend/handlers/auth.go`
- Understand `LoginHandler()` step by step
- Understand `SignupHandler()`
- Understand `HasuraLoginHandler()` and `HasuraSignupHandler()`
- Trace through the code with a pen

#### 3. Study JWT (15 min)
- Read `backend/utils/jwt.go`
- Understand how JWT is generated
- Understand claims structure
- Understand Hasura claims

### Afternoon (1.5 hours)

#### 4. Study Auth Middleware (30 min)
- Find `AuthMiddleware` function
- Understand how it verifies JWT
- Understand how it extracts user_id
- Understand context

#### 5. Study Recipe Handlers (45 min)
- Read `backend/handlers/recipes.go` or `recipes_handler.go`
- Understand `CreateRecipeHandler()`
- Understand transactions
- Understand why transactions are needed
- Study `GetRecipesHandler()`
- Study `GetRecipeByIDHandler()`

#### 6. Practice (15 min)
- Explain authentication flow to yourself
- Explain how a recipe is created
- Explain why transactions are used

### ✅ Day 3 Checklist
- [ ] I understand authentication flow
- [ ] I can explain how JWT works
- [ ] I understand middleware pattern
- [ ] I can explain how recipes are created
- [ ] I understand database transactions
- [ ] I can trace a request from start to finish

---

## 📆 Day 4: Frontend Foundation (3 hours)

### Morning (1.5 hours)

#### 1. Read `COMPLETE_PROJECT_GUIDE.md` - Day 4 Section (30 min)
- Understand Vue.js basics
- Understand Nuxt.js structure
- Understand file-based routing

#### 2. Study Frontend Structure (30 min)
- Explore `frontend/nuxt-app/` folder
- Understand `pages/` directory
- Understand `layouts/` directory
- Understand `composables/` directory
- Understand `plugins/` directory

#### 3. Study Login Page (30 min)
- Read `frontend/nuxt-app/pages/login.vue`
- Understand template syntax
- Understand `v-model`
- Understand `@submit.prevent`
- Understand `ref()` and reactive data
- Understand `fetch()` API calls

### Afternoon (1.5 hours)

#### 4. Study useAuth Composable (30 min)
- Read `frontend/nuxt-app/composables/useAuth.js`
- Understand what composables are
- Understand `useState()`
- Understand `computed()`
- Understand how login/logout works

#### 5. Study Apollo Client (30 min)
- Read `frontend/nuxt-app/plugins/apollo.client.ts`
- Understand GraphQL client setup
- Understand why we use GraphQL
- Understand difference from REST

#### 6. Study Recipe Pages (30 min)
- Read `frontend/nuxt-app/pages/create.vue`
- Read `frontend/nuxt-app/pages/recipes/[id].vue`
- Understand dynamic routes
- Understand how data is fetched
- Understand how forms work

### ✅ Day 4 Checklist
- [ ] I understand Vue.js component structure
- [ ] I understand Nuxt.js file-based routing
- [ ] I understand reactive data (`ref`, `computed`)
- [ ] I understand how frontend calls backend
- [ ] I understand GraphQL vs REST
- [ ] I can explain how login works on frontend

---

## 📆 Day 5: Integration & Practice (3 hours)

### Morning (1.5 hours)

#### 1. Read `COMPLETE_PROJECT_GUIDE.md` - Day 5 Section (30 min)
- Understand file uploads
- Understand payment integration
- Understand Hasura Actions
- Understand event triggers

#### 2. Study File Upload (20 min)
- Read `backend/handlers/extra.go` → `UploadFileHandler()`
- Understand multipart form
- Understand file saving
- Understand URL generation

#### 3. Study Payment (20 min)
- Read `backend/handlers/payment.go`
- Understand payment flow
- Understand webhook callback
- Understand Chapa integration

#### 4. Study Hasura Integration (20 min)
- Read `backend/handlers/auth.go` → Hasura handlers
- Read `backend/handlers/events.go`
- Understand Actions vs Events
- Understand why we need them

### Afternoon (1.5 hours)

#### 5. Study Social Features (30 min)
- Read `backend/handlers/social.go`
- Understand likes, comments, ratings
- Understand toggle operations
- Understand many-to-many relationships

#### 6. Practice Presentation (60 min)
- Use `QUICK_REFERENCE_CHEATSHEET.md`
- Practice explaining:
  - Project overview (30 sec)
  - Architecture (2 min)
  - Authentication flow (2 min)
  - Database structure (2 min)
  - API endpoints (2 min)
  - One complete feature (login or create recipe) (5 min)
- Record yourself or practice with a friend
- Time yourself

### ✅ Day 5 Checklist
- [ ] I understand file uploads
- [ ] I understand payment flow
- [ ] I understand Hasura Actions and Events
- [ ] I understand social features
- [ ] I can explain the complete project in 15 minutes
- [ ] I feel confident answering questions

---

## 🎯 Final Preparation (Day Before Presentation)

### Review Session (2 hours)

1. **Quick Review** (30 min)
   - Read `QUICK_REFERENCE_CHEATSHEET.md`
   - Review key concepts

2. **Code Walkthrough** (45 min)
   - Open key files
   - Trace through one complete flow:
     - User logs in → Creates recipe → Uploads image → Likes recipe
   - Explain each step

3. **Practice Q&A** (30 min)
   - Review "Common Questions" from cheat sheet
   - Practice answers
   - Think of additional questions boss might ask

4. **Final Practice** (15 min)
   - Give full presentation to yourself
   - Time it (should be 10-15 minutes)
   - Make sure you cover:
     - What the project does
     - Architecture
     - Key features
     - Technology choices
     - One detailed walkthrough

---

## 📝 Study Tips

### While Studying:

1. **Take Notes**
   - Write down key concepts
   - Draw diagrams
   - List questions

2. **Code Along**
   - Open files while reading
   - Trace through code with your finger
   - Add comments to understand

3. **Explain Out Loud**
   - Explain concepts to yourself
   - Pretend you're teaching someone
   - This helps retention

4. **Draw Diagrams**
   - Architecture diagram
   - Database relationships
   - Request flow
   - Authentication flow

5. **Practice Questions**
   - "How does X work?"
   - "Why did we use Y?"
   - "What happens when Z?"

### Common Mistakes to Avoid:

1. ❌ Just reading without understanding
2. ❌ Not looking at actual code
3. ❌ Memorizing without understanding
4. ❌ Not practicing explanation
5. ❌ Not preparing for questions

### What to Focus On:

1. ✅ **Big Picture First**: Understand architecture before details
2. ✅ **Flow**: Understand how data flows through the system
3. ✅ **Why**: Understand why each technology was chosen
4. ✅ **One Feature Deeply**: Master one feature completely (like login)
5. ✅ **Connections**: Understand how frontend connects to backend

---

## 🎤 Presentation Structure

### Suggested 15-Minute Presentation:

1. **Introduction** (1 min)
   - "This is a Food Recipes App I built..."
   - Show the app running

2. **Architecture Overview** (2 min)
   - Draw diagram: Browser → Frontend → Backend → Database
   - Explain tech stack
   - Explain why each technology

3. **Key Features** (2 min)
   - Authentication
   - Recipe management
   - Social features
   - Payments

4. **Deep Dive: Login Flow** (5 min)
   - Show login page
   - Explain frontend code
   - Explain backend code
   - Show database query
   - Show JWT generation
   - Trace complete flow

5. **Database Structure** (2 min)
   - Show main tables
   - Explain relationships
   - Explain why this structure

6. **Questions & Answers** (3 min)
   - Be ready for questions
   - Use cheat sheet

---

## 💪 Confidence Building

### Remember:

1. **You Built This**: Even with help, you understand it now
2. **You Studied**: You've spent time learning
3. **You Can Explain**: You can walk through the code
4. **It's Okay to Say "I Don't Know"**: But follow with "Let me check the code"

### Before Presentation:

- ✅ Get good sleep
- ✅ Review cheat sheet
- ✅ Have code open and ready
- ✅ Take deep breaths
- ✅ Remember: You've got this!

---

## 🎯 Success Criteria

You're ready when you can:

- [ ] Explain what the project does in 30 seconds
- [ ] Draw the architecture diagram
- [ ] Explain authentication flow
- [ ] Walk through one complete feature (login or create recipe)
- [ ] Explain database structure
- [ ] Answer "why" questions about technology choices
- [ ] Trace a request from frontend to database and back
- [ ] Explain security measures
- [ ] Feel confident (even if nervous)

---

**Good luck! You've got this! 🚀**

Remember: Understanding is more important than memorization. If you understand the concepts, you can explain them in your own words.



