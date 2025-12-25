# Clean existing data and reimport
$dbUrl = "postgresql://postgres:IUQsoQlDjxLLyERgGppSbRnOPdmOZUbu@ballast.proxy.rlwy.net:51054/railway"

Write-Host "`nStep 1: Cleaning existing data..." -ForegroundColor Cyan

$cleanSql = @"
-- Disable foreign key checks temporarily
SET session_replication_role = 'replica';

-- Truncate all data tables (preserve structure)
TRUNCATE TABLE 
    users, 
    categories, 
    recipes, 
    recipe_ingredients, 
    recipe_steps, 
    recipe_images,
    bookmarks,
    comments,
    likes,
    purchases,
    ratings
CASCADE;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE categories_id_seq RESTART WITH 1;
ALTER SEQUENCE recipes_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_ingredients_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_steps_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_images_id_seq RESTART WITH 1;
ALTER SEQUENCE comments_id_seq RESTART WITH 1;
ALTER SEQUENCE purchases_id_seq RESTART WITH 1;

-- Re-enable foreign key checks
SET session_replication_role = 'origin';
"@

$cleanSql | docker run -i --rm postgres:15 psql $dbUrl

if ($LASTEXITCODE -ne 0) {
    Write-Host "Clean failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Data cleaned!" -ForegroundColor Green
Write-Host "`nStep 2: Importing fresh data..." -ForegroundColor Cyan

# Extract only COPY statements from backup (data only)
$dataOnly = Get-Content database_backup_clean.sql | Select-String -Pattern "^COPY " -Context 0,1000

# For now, let's just reimport the full backup - it should work now that tables are empty
Get-Content database_backup_clean.sql | docker run -i --rm postgres:15 psql $dbUrl 2>&1 | Select-String -Pattern "COPY|ERROR" 

Write-Host "`nImport complete! Checking results..." -ForegroundColor Green

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
UNION ALL SELECT 'recipe_steps', COUNT(*) FROM recipe_steps;
'@
docker run -i --rm postgres:15 psql $dbUrl -c $checkSql

