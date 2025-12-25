# Fix Image URLs in Railway Database
# Replace localhost URLs with Railway backend URL

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  FIX IMAGE URLS IN RAILWAY DATABASE" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Step 1: Get your Railway Backend URL" -ForegroundColor Yellow
Write-Host "1. Go to Railway → Backend service → Settings → Domain" -ForegroundColor White
Write-Host "2. Copy the URL (e.g., https://backend-production-xxxx.up.railway.app)`n" -ForegroundColor White

$railwayBackendUrl = Read-Host "Paste Railway Backend URL here (without trailing slash)"

if ([string]::IsNullOrWhiteSpace($railwayBackendUrl)) {
    Write-Host "Backend URL is required!" -ForegroundColor Red
    exit 1
}

# Remove trailing slash
$railwayBackendUrl = $railwayBackendUrl.TrimEnd('/')

$dbUrl = "postgresql://postgres:IUQsoQlDjxLLyERgGppSbRnOPdmOZUbu@ballast.proxy.rlwy.net:51054/railway"

Write-Host "`nStep 2: Updating image URLs..." -ForegroundColor Yellow
Write-Host "Replacing: http://localhost:8081" -ForegroundColor White
Write-Host "With: $railwayBackendUrl`n" -ForegroundColor White

# Update recipe_images table
$updateImages = @"
UPDATE recipe_images 
SET url = REPLACE(url, 'http://localhost:8081', '$railwayBackendUrl')
WHERE url LIKE 'http://localhost:8081%';
"@

Write-Host "Updating recipe_images..." -NoNewline
$result1 = $updateImages | docker run -i --rm postgres:15 psql $dbUrl -c $updateImages 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host " Done!" -ForegroundColor Green
} else {
    Write-Host " Failed!" -ForegroundColor Red
    Write-Host $result1
}

# Update recipes thumbnail_url
$updateThumbnails = @"
UPDATE recipes 
SET thumbnail_url = REPLACE(thumbnail_url, 'http://localhost:8081', '$railwayBackendUrl')
WHERE thumbnail_url LIKE 'http://localhost:8081%';
"@

Write-Host "Updating recipes thumbnail_url..." -NoNewline
$result2 = $updateThumbnails | docker run -i --rm postgres:15 psql $dbUrl -c $updateThumbnails 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host " Done!" -ForegroundColor Green
} else {
    Write-Host " Failed!" -ForegroundColor Red
    Write-Host $result2
}

# Update recipe_steps image_url
$updateSteps = @"
UPDATE recipe_steps 
SET image_url = REPLACE(image_url, 'http://localhost:8081', '$railwayBackendUrl')
WHERE image_url LIKE 'http://localhost:8081%';
"@

Write-Host "Updating recipe_steps image_url..." -NoNewline
$result3 = $updateSteps | docker run -i --rm postgres:15 psql $dbUrl -c $updateSteps 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host " Done!" -ForegroundColor Green
} else {
    Write-Host " Failed!" -ForegroundColor Red
    Write-Host $result3
}

Write-Host "`nStep 3: Verifying updates..." -ForegroundColor Yellow

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

$remaining = docker run -i --rm postgres:15 psql $dbUrl -c $checkSql

Write-Host "`nRemaining localhost URLs:" -ForegroundColor Cyan
$remaining

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  UPDATE COMPLETE!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Add API_URL variable to Railway Backend:" -ForegroundColor White
Write-Host "   Name: API_URL" -ForegroundColor Green
Write-Host "   Value: $railwayBackendUrl" -ForegroundColor Green
Write-Host "2. Redeploy backend (Railway will auto-redeploy)" -ForegroundColor White
Write-Host "3. Test images in your app!" -ForegroundColor White

