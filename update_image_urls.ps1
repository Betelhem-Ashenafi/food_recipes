# Update Image URLs in Railway Database
param(
    [Parameter(Mandatory=$true)]
    [string]$BackendUrl
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  UPDATING IMAGE URLS" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Remove trailing slash
$BackendUrl = $BackendUrl.TrimEnd('/')

$dbUrl = "postgresql://postgres:IUQsoQlDjxLLyERgGppSbRnOPdmOZUbu@ballast.proxy.rlwy.net:51054/railway"

Write-Host "Replacing: http://localhost:8081" -ForegroundColor Yellow
Write-Host "With: $BackendUrl`n" -ForegroundColor Green

# Update recipe_images table
Write-Host "Updating recipe_images..." -NoNewline
$updateImages = "UPDATE recipe_images SET url = REPLACE(url, 'http://localhost:8081', '$BackendUrl') WHERE url LIKE 'http://localhost:8081%';"
$result1 = echo $updateImages | docker run -i --rm postgres:15 psql $dbUrl 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host " Done!" -ForegroundColor Green
} else {
    Write-Host " Failed!" -ForegroundColor Red
}

# Update recipes thumbnail_url
Write-Host "Updating recipes thumbnail_url..." -NoNewline
$updateThumbnails = "UPDATE recipes SET thumbnail_url = REPLACE(thumbnail_url, 'http://localhost:8081', '$BackendUrl') WHERE thumbnail_url LIKE 'http://localhost:8081%';"
$result2 = echo $updateThumbnails | docker run -i --rm postgres:15 psql $dbUrl 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host " Done!" -ForegroundColor Green
} else {
    Write-Host " Failed!" -ForegroundColor Red
}

# Update recipe_steps image_url
Write-Host "Updating recipe_steps image_url..." -NoNewline
$updateSteps = "UPDATE recipe_steps SET image_url = REPLACE(image_url, 'http://localhost:8081', '$BackendUrl') WHERE image_url LIKE 'http://localhost:8081%';"
$result3 = echo $updateSteps | docker run -i --rm postgres:15 psql $dbUrl 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host " Done!" -ForegroundColor Green
} else {
    Write-Host " Failed!" -ForegroundColor Red
}

Write-Host "`nVerifying updates..." -ForegroundColor Yellow

# Check results
$checkSql = @'
SELECT 
    'recipe_images' as table_name, COUNT(*) as localhost_count 
FROM recipe_images 
WHERE url LIKE 'http://localhost:8081%'
UNION ALL
SELECT 
    'recipes', COUNT(*) 
FROM recipes 
WHERE thumbnail_url LIKE 'http://localhost:8081%'
UNION ALL
SELECT 
    'recipe_steps', COUNT(*) 
FROM recipe_steps 
WHERE image_url LIKE 'http://localhost:8081%';
'@

$remaining = echo $checkSql | docker run -i --rm postgres:15 psql $dbUrl -c $checkSql

Write-Host "`nRemaining localhost URLs:" -ForegroundColor Cyan
$remaining

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  UPDATE COMPLETE!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Cyan

