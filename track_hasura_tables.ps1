# Track All Tables in Hasura (Railway)
# This script automatically tracks all tables in your Hasura instance

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  TRACK TABLES IN HASURA" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Step 1: Get your Hasura URL" -ForegroundColor Yellow
Write-Host "1. Go to Railway → Hasura service → Settings → Domain" -ForegroundColor White
Write-Host "2. Copy the URL (e.g., https://captivating-wholeness-production.up.railway.app)`n" -ForegroundColor White

$hasuraUrl = Read-Host "Paste Hasura URL here (without /v1/graphql)"

if ([string]::IsNullOrWhiteSpace($hasuraUrl)) {
    Write-Host "Hasura URL is required!" -ForegroundColor Red
    exit 1
}

# Remove trailing slash
$hasuraUrl = $hasuraUrl.TrimEnd('/')

$adminSecret = "myhasurasecret"

$headers = @{
    "Content-Type" = "application/json"
    "X-Hasura-Admin-Secret" = $adminSecret
}

$tables = @(
    "users",
    "categories", 
    "recipes", 
    "recipe_ingredients", 
    "recipe_steps", 
    "recipe_images",
    "likes",
    "bookmarks",
    "comments",
    "ratings",
    "purchases"
)

Write-Host "`nStep 2: Tracking tables..." -ForegroundColor Yellow
Write-Host "Connecting to: $hasuraUrl`n" -ForegroundColor White

$successCount = 0
$failCount = 0

foreach ($table in $tables) {
    $payload = @{
        type = "pg_track_table"
        args = @{
            source = "default"
            table = @{
                schema = "public"
                name = $table
            }
        }
    } | ConvertTo-Json -Depth 10

    Write-Host "Tracking table: '$table' ... " -NoNewline
    try {
        $response = Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $payload
        Write-Host "Success!" -ForegroundColor Green
        $successCount++
    } catch {
        $errorMessage = $_.Exception.Message
        if ($_.Exception.Response) {
            $stream = $_.Exception.Response.GetResponseStream()
            $reader = [System.IO.StreamReader]::new($stream)
            $errBody = $reader.ReadToEnd()
            
            if ($errBody -match "already tracked" -or $errBody -match "already exists") {
                Write-Host "Already tracked!" -ForegroundColor Yellow
                $successCount++
            } else {
                Write-Host "Failed!" -ForegroundColor Red
                Write-Host "  Error: $errBody" -ForegroundColor Red
                $failCount++
            }
        } else {
            Write-Host "Failed!" -ForegroundColor Red
            Write-Host "  Error: $errorMessage" -ForegroundColor Red
            $failCount++
        }
    }
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  RESULTS" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Successfully tracked: $successCount tables" -ForegroundColor Green
if ($failCount -gt 0) {
    Write-Host "Failed: $failCount tables" -ForegroundColor Red
}

Write-Host "`nNext steps:" -ForegroundColor Cyan
Write-Host "1. Go to Hasura Console: $hasuraUrl/console" -ForegroundColor White
Write-Host "2. Enter admin secret: $adminSecret" -ForegroundColor White
Write-Host "3. Click 'Data' tab → You should see all tables!" -ForegroundColor White

