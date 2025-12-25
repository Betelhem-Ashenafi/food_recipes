# Export Local Database and Import to Railway
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  EXPORT LOCAL → IMPORT RAILWAY" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$railwayDbUrl = "postgresql://postgres:IUQsoQlDjxLLyERgGppSbRnOPdmOZUbu@ballast.proxy.rlwy.net:51054/railway"
$backupFile = "local_database_backup.sql"

Write-Host "Step 1: Exporting local database..." -ForegroundColor Yellow
Write-Host "This may take 30 seconds...`n" -ForegroundColor White

# Export local database
docker exec docker-postgres-1 pg_dump -U fooduser -d foodrecipes --data-only --column-inserts > $backupFile

if ($LASTEXITCODE -ne 0) {
    Write-Host "Export failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Export complete! File: $backupFile" -ForegroundColor Green
Write-Host "File size: $((Get-Item $backupFile).Length / 1KB) KB`n" -ForegroundColor White

Write-Host "Step 2: Cleaning Railway database..." -ForegroundColor Yellow

# Clean Railway database (truncate all tables)
$cleanSql = @"
SET session_replication_role = 'replica';
TRUNCATE TABLE 
    users, categories, recipes, recipe_ingredients, recipe_steps, 
    recipe_images, bookmarks, comments, likes, purchases, ratings
CASCADE;
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE categories_id_seq RESTART WITH 1;
ALTER SEQUENCE recipes_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_ingredients_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_steps_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_images_id_seq RESTART WITH 1;
ALTER SEQUENCE comments_id_seq RESTART WITH 1;
ALTER SEQUENCE purchases_id_seq RESTART WITH 1;
SET session_replication_role = 'origin';
"@

$cleanSql | docker run -i --rm postgres:15 psql $railwayDbUrl

if ($LASTEXITCODE -ne 0) {
    Write-Host "Clean failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Railway database cleaned!`n" -ForegroundColor Green

Write-Host "Step 3: Importing data to Railway..." -ForegroundColor Yellow
Write-Host "This may take 1-2 minutes...`n" -ForegroundColor White

# Import to Railway
Get-Content $backupFile | docker run -i --rm postgres:15 psql $railwayDbUrl 2>&1 | Select-String -Pattern "INSERT|ERROR" | Select-Object -First 20

Write-Host "`nStep 4: Verifying import..." -ForegroundColor Yellow

# Check what was imported
$checkSql = @'
SELECT 'users' as table_name, COUNT(*) as count FROM users 
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

docker run -i --rm postgres:15 psql $railwayDbUrl -c $checkSql

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  IMPORT COMPLETE!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Cyan

