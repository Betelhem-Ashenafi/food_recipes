# Export Local Database to Railway
# This will export your existing local database with all your data

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  EXPORT LOCAL DATABASE" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Local database connection
$localDb = "postgres://fooduser:foodpass@localhost:5433/foodrecipes"

# Export database
Write-Host "[1/2] Exporting local database..." -ForegroundColor Yellow
$dumpFile = "database_backup.sql"

# Check if pg_dump is available
$pgDumpPath = Get-Command pg_dump -ErrorAction SilentlyContinue
if (-not $pgDumpPath) {
    Write-Host "❌ pg_dump not found. Please install PostgreSQL client tools." -ForegroundColor Red
    Write-Host "   Or use Docker to export:" -ForegroundColor Yellow
    Write-Host "   docker exec -t postgres pg_dump -U fooduser foodrecipes > database_backup.sql" -ForegroundColor Cyan
    exit 1
}

# Export using pg_dump
pg_dump "$localDb" -F p -f $dumpFile

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Database exported to: $dumpFile" -ForegroundColor Green
    Write-Host "`nFile size: $((Get-Item $dumpFile).Length / 1KB) KB" -ForegroundColor Cyan
} else {
    Write-Host "❌ Export failed" -ForegroundColor Red
    exit 1
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  NEXT STEPS:" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan
Write-Host "1. Go to Railway → Postgres service → Connect" -ForegroundColor White
Write-Host "2. Copy the DATABASE_URL" -ForegroundColor White
Write-Host "3. Run this command:" -ForegroundColor White
Write-Host "   psql `"YOUR_RAILWAY_DATABASE_URL`" < database_backup.sql" -ForegroundColor Cyan
Write-Host "`nOR use Hasura Console:" -ForegroundColor Yellow
Write-Host "1. Go to: https://captivating-wholeness-production.up.railway.app" -ForegroundColor White
Write-Host "2. Data → SQL → Paste the SQL from database_backup.sql" -ForegroundColor White
Write-Host "3. Click Run" -ForegroundColor White

