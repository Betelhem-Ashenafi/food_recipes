# Import Database to Railway - PowerShell Compatible
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  IMPORT DATABASE TO RAILWAY" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Step 1: Get Railway DATABASE_URL" -ForegroundColor Yellow
Write-Host "1. Go to Railway → Postgres service → Variables tab" -ForegroundColor White
Write-Host "2. Copy the DATABASE_URL value`n" -ForegroundColor White

$railwayDbUrl = Read-Host "Paste Railway DATABASE_URL here"

if ([string]::IsNullOrWhiteSpace($railwayDbUrl)) {
    Write-Host "❌ DATABASE_URL is required!" -ForegroundColor Red
    exit 1
}

$backupFile = "database_backup_clean.sql"
if (-not (Test-Path $backupFile)) {
    Write-Host "❌ $backupFile not found!" -ForegroundColor Red
    exit 1
}

Write-Host "`nStep 2: Importing database..." -ForegroundColor Yellow
Write-Host "This may take 1-2 minutes...`n" -ForegroundColor White

# Use Get-Content and pipe to Docker (PowerShell compatible)
Get-Content $backupFile | docker run -i --rm postgres:15 psql $railwayDbUrl

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ Database imported successfully!" -ForegroundColor Green
    Write-Host "`nNext steps:" -ForegroundColor Cyan
    Write-Host "1. Go to Hasura Console: https://captivating-wholeness-production.up.railway.app" -ForegroundColor White
    Write-Host "2. Click 'Data' tab → 'Track All'" -ForegroundColor White
    Write-Host "3. Update action URLs in 'Actions' tab" -ForegroundColor White
} else {
    Write-Host "`n❌ Import failed. Check the error above." -ForegroundColor Red
}

