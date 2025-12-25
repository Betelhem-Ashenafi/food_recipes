# Copy All Data from Local Database to Railway
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  COPY LOCAL DATABASE → RAILWAY" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$railwayDbUrl = "postgresql://postgres:IUQsoQlDjxLLyERgGppSbRnOPdmOZUbu@ballast.proxy.rlwy.net:51054/railway"
$backupFile = "local_to_railway_backup.sql"

Write-Host "Step 1: Exporting local database (with data)..." -ForegroundColor Yellow
Write-Host "This may take 30 seconds...`n" -ForegroundColor White

# Export local database with data only (no schema)
docker exec docker-postgres-1 pg_dump -U fooduser -d foodrecipes `
    --data-only `
    --column-inserts `
    --no-owner `
    --no-privileges `
    --exclude-table=hdb_* `
    --exclude-table=event_* `
    > $backupFile

if ($LASTEXITCODE -ne 0) {
    Write-Host "Export failed!" -ForegroundColor Red
    exit 1
}

$fileSize = (Get-Item $backupFile).Length / 1KB
Write-Host "Export complete! File: $backupFile ($([math]::Round($fileSize, 2)) KB)`n" -ForegroundColor Green

Write-Host "Step 2: Cleaning Railway database..." -ForegroundColor Yellow

# Clean Railway database (truncate all tables in correct order)
$cleanSql = @"
-- Disable foreign key checks
SET session_replication_role = 'replica';

-- Truncate in reverse dependency order
TRUNCATE TABLE 
    recipe_steps,
    recipe_ingredients,
    recipe_images,
    comments,
    likes,
    bookmarks,
    purchases,
    ratings,
    recipes,
    categories,
    users
CASCADE;

-- Reset sequences
ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS categories_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS recipes_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS recipe_ingredients_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS recipe_steps_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS recipe_images_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS comments_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS purchases_id_seq RESTART WITH 1;

-- Re-enable foreign key checks
SET session_replication_role = 'origin';
"@

$cleanSql | docker run -i --rm postgres:15 psql $railwayDbUrl | Out-Null

if ($LASTEXITCODE -ne 0) {
    Write-Host "Clean failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Railway database cleaned!`n" -ForegroundColor Green

Write-Host "Step 3: Importing data to Railway..." -ForegroundColor Yellow
Write-Host "This may take 1-2 minutes...`n" -ForegroundColor White

# Import to Railway (suppress most output, only show errors)
$importOutput = Get-Content $backupFile | docker run -i --rm postgres:15 psql $railwayDbUrl 2>&1

# Show only errors
$errors = $importOutput | Select-String -Pattern "ERROR"
if ($errors) {
    Write-Host "Import completed with some errors:" -ForegroundColor Yellow
    $errors | Select-Object -First 10
} else {
    Write-Host "Import completed successfully!`n" -ForegroundColor Green
}

Write-Host "`nStep 4: Verifying import..." -ForegroundColor Yellow

# Check what was imported
$checkSql = @'
SELECT 
    'users' as table_name, COUNT(*) as count FROM users 
UNION ALL SELECT 'recipes', COUNT(*) FROM recipes 
UNION ALL SELECT 'categories', COUNT(*) FROM categories 
UNION ALL SELECT 'bookmarks', COUNT(*) FROM bookmarks 
UNION ALL SELECT 'comments', COUNT(*) FROM comments 
UNION ALL SELECT 'likes', COUNT(*) FROM likes 
UNION ALL SELECT 'recipe_images', COUNT(*) FROM recipe_images 
UNION ALL SELECT 'recipe_ingredients', COUNT(*) FROM recipe_ingredients 
UNION ALL SELECT 'recipe_steps', COUNT(*) FROM recipe_steps 
UNION ALL SELECT 'purchases', COUNT(*) FROM purchases 
UNION ALL SELECT 'ratings', COUNT(*) FROM ratings 
ORDER BY table_name;
'@

Write-Host "`nData in Railway database:" -ForegroundColor Cyan
docker run -i --rm postgres:15 psql $railwayDbUrl -c $checkSql

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  COPY COMPLETE!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Cyan
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Connect Railway Hasura to Railway database" -ForegroundColor White
Write-Host "2. Track tables in Hasura Console" -ForegroundColor White
Write-Host "3. Test your app!" -ForegroundColor White

