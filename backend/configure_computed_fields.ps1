$headers = @{
    "Content-Type" = "application/json"
    "X-Hasura-Admin-Secret" = "myhasurasecret"
}

$computedFields = @(
    @{
        table = "recipes"
        name = "average_rating"
        function = "recipe_average_rating"
    },
    @{
        table = "recipes"
        name = "likes_count"
        function = "recipe_likes_count"
    }
)

Write-Host "Configuring Computed Fields in Hasura..."

foreach ($field in $computedFields) {
    $payload = @{
        type = "pg_create_computed_field"
        args = @{
            source = "default"
            table = @{
                schema = "public"
                name = $field.table
            }
            name = $field.name
            definition = @{
                function = @{
                    schema = "public"
                    name = $field.function
                }
            }
        }
    } | ConvertTo-Json -Depth 10

    Write-Host "Adding computed field '$($field.name)' to '$($field.table)'..." -NoNewline
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $payload
        Write-Host " Success." -ForegroundColor Green
    } catch {
        $stream = $_.Exception.Response.GetResponseStream()
        $reader = [System.IO.StreamReader]::new($stream)
        $errBody = $reader.ReadToEnd()
        
        if ($errBody -match "already exists") {
            Write-Host " Already exists." -ForegroundColor Yellow
        } else {
            Write-Host " Failed." -ForegroundColor Red
            Write-Host "Error details: $errBody"
        }
    }
}
