# Full Stack Startup Script
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  FOOD RECIPES APP - FULL STACK STARTUP" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# 1. Start Docker Services
Write-Host "[1/4] Starting Docker services (Postgres + Hasura)..." -ForegroundColor Yellow
Set-Location docker
docker-compose up -d
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Docker services started" -ForegroundColor Green
} else {
    Write-Host "❌ Failed to start Docker services" -ForegroundColor Red
    exit 1
}

# Wait for Hasura to be ready
Write-Host "`n[2/4] Waiting for Hasura to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10
Write-Host "✅ Hasura should be ready" -ForegroundColor Green

# 2. Start Backend
Write-Host "`n[3/4] Starting Backend (Golang)..." -ForegroundColor Yellow
Set-Location ../backend
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PWD'; Write-Host 'Starting Backend...' -ForegroundColor Green; go run main.go"
Start-Sleep -Seconds 3
Write-Host "✅ Backend started on http://localhost:8081" -ForegroundColor Green

# 3. Start Frontend
Write-Host "`n[4/4] Starting Frontend (Nuxt)..." -ForegroundColor Yellow
Set-Location ../frontend/nuxt-app
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PWD'; Write-Host 'Starting Frontend...' -ForegroundColor Green; npm run dev"
Start-Sleep -Seconds 5
Write-Host "✅ Frontend started on http://localhost:3000" -ForegroundColor Green

# 4. Display Status
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  FULL STACK IS RUNNING!" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Services:" -ForegroundColor White
Write-Host "  ✅ Postgres      - http://localhost:5433" -ForegroundColor Green
Write-Host "  ✅ Hasura        - http://localhost:8080" -ForegroundColor Green
Write-Host "  ✅ Backend (Go)  - http://localhost:8081" -ForegroundColor Green
Write-Host "  ✅ Frontend      - http://localhost:3000" -ForegroundColor Green

Write-Host "`nAccess Your App:" -ForegroundColor White
Write-Host "  🌐 Frontend      - http://localhost:3000" -ForegroundColor Cyan
Write-Host "  🔐 Login         - http://localhost:3000/login" -ForegroundColor Cyan
Write-Host "  🏠 Home          - http://localhost:3000/home" -ForegroundColor Cyan
Write-Host "  ➕ Create Recipe - http://localhost:3000/create" -ForegroundColor Cyan
Write-Host "  ⚙️  Hasura Console - http://localhost:8080 (secret: myhasurasecret)" -ForegroundColor Cyan

Write-Host "`nPress Ctrl+C to stop this script (services will continue running in background)" -ForegroundColor Yellow
Write-Host "`nTo stop all services:" -ForegroundColor Yellow
Write-Host "  cd docker; docker-compose down" -ForegroundColor Gray

# Keep script running
Read-Host "`nPress Enter to exit"












