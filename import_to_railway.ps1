# Import Database to Railway - No Chunks Needed!
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  IMPORT DATABASE TO RAILWAY" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Step 1: Get Railway DATABASE_URL" -ForegroundColor Yellow
Write-Host "1. Go to Railway → Postgres service → Variables tab" -ForegroundColor White
Write-Host "2. Copy the DATABASE_URL value" -ForegroundColor White
Write-Host "3. Paste it below when prompted`n" -ForegroundColor White

$railwayDbUrl = Read-Host "Paste Railway DATABASE_URL here"

if ([string]::IsNullOrWhiteSpace($railwayDbUrl)) {
    Write-Host "❌ DATABASE_URL is required!" -ForegroundColor Red
    exit 1
}

Write-Host "`nStep 2: Importing database..." -ForegroundColor Yellow
Write-Host "This may take 1-2 minutes...`n" -ForegroundColor White

# Check if psql is available
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psqlPath) {
    Write-Host "❌ psql not found. Installing..." -ForegroundColor Yellow
    Write-Host "Please install PostgreSQL client tools first:" -ForegroundColor Yellow
    Write-Host "Download from: https://www.postgresql.org/download/windows/" -ForegroundColor Cyan
    Write-Host "`nOR use Docker method below:`n" -ForegroundColor Yellow
    
    Write-Host "Alternative: Use Docker to import" -ForegroundColor Cyan
    Write-Host "docker run -i --rm postgres:15 psql `"$railwayDbUrl`" < database_backup.sql" -ForegroundColor Green
    exit 1
}

# Import using psql
$backupFile = "database_backup.sql"
if (-not (Test-Path $backupFile)) {
    Write-Host "❌ database_backup.sql not found!" -ForegroundColor Red
    exit 1
}

Write-Host "Importing $backupFile to Railway..." -ForegroundColor Yellow
Get-Content $backupFile | & psql $railwayDbUrl

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ Database imported successfully!" -ForegroundColor Green
    Write-Host "`nNext steps:" -ForegroundColor Cyan
    Write-Host "1. Go to Hasura Console: https://captivating-wholeness-production.up.railway.app" -ForegroundColor White
    Write-Host "2. Click 'Data' tab → 'Track All'" -ForegroundColor White
    Write-Host "3. Update action URLs in 'Actions' tab" -ForegroundColor White
} else {
    Write-Host "`n❌ Import failed. Check the error above." -ForegroundColor Red
}

