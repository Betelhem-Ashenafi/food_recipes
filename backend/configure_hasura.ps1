$headers = @{
    "Content-Type" = "application/json"
    "X-Hasura-Admin-Secret" = "myhasurasecret"
}

$tables = @("users", "categories", "recipes", "recipe_ingredients", "recipe_steps", "recipe_images", "likes", "bookmarks", "comments", "ratings")

Write-Host "Connecting to Hasura at http://localhost:8080..."

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

    Write-Host "Tracking table: '$table' ..." -NoNewline
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $payload
        Write-Host " Success." -ForegroundColor Green
    } catch {
        # Check if it's an "already tracked" error
        $stream = $_.Exception.Response.GetResponseStream()
        $reader = [System.IO.StreamReader]::new($stream)
        $errBody = $reader.ReadToEnd()
        
        if ($errBody -match "already tracked") {
            Write-Host " Already tracked." -ForegroundColor Yellow
        } else {
            Write-Host " Failed." -ForegroundColor Red
            Write-Host "Error details: $errBody"
        }
    }
}

Write-Host "`nConfiguration Complete! Your API is ready." -ForegroundColor Cyan