# Food Recipes App

A modern food recipe website where users can browse, share, and interact with recipes. Built with Golang, Nuxt 4, Vue 3, Hasura, Postgres, and Docker.

## Features
- Browse recipes by category, creator, time, ingredients, and title
- Signup/login, create/edit/delete recipes
- Upload multiple images, select featured image
- Dynamic steps and ingredients
- Like, bookmark, comment, rate, and buy recipes
- JWT authentication
- Hasura events, actions, permissions, triggers, computed/generated columns
- Chapa integration for payments
- Attractive, responsive UI

## Tech Stack
- Backend: Golang >v1.22, GraphQL
- Frontend: Nuxt 4, Vue 3, Vite, TailwindCSS, Urql
- Database: Postgres
- API Gateway: Hasura
- Docker for local development

## Setup & Run

1. **Start Infrastructure (Hasura & Postgres)**
   ```bash
   cd docker
   docker-compose up -d
   ```

2. **Start Backend**
   ```bash
   cd backend
   go run main.go
   # Or build and run
   go build -o backend.exe
   ./backend.exe
   ```
   Backend runs on `http://localhost:8081`.

3. **Start Frontend**
   ```bash
   cd frontend/nuxt-app
   npm install
   npm run dev
   ```
   Frontend runs on `http://localhost:3000`.

## Notes
- Ensure `CHAPA_SECRET_KEY` is set in environment or `backend/handlers/payment.go` (currently hardcoded for dev).
- Hasura Console: `http://localhost:8080` (Secret: `myhasurasecret`).
