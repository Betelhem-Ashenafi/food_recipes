# Configure Hasura Relationships
$headers = @{
    "Content-Type" = "application/json"
    "X-Hasura-Admin-Secret" = "myhasurasecret"
}

Write-Host "Configuring Hasura Relationships..." -ForegroundColor Cyan

# Relationship: recipes -> users (user_id -> users.id)
$recipesToUsers = @{
    type = "pg_create_object_relationship"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "recipes"
        }
        name = "user"
        using = @{
            foreign_key_constraint_on = "user_id"
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Creating recipes -> users relationship..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $recipesToUsers
    Write-Host " Success." -ForegroundColor Green
} catch {
    if ($_.Exception.Response) {
        $stream = $_.Exception.Response.GetResponseStream()
        $reader = [System.IO.StreamReader]::new($stream)
        $errBody = $reader.ReadToEnd()
        if ($errBody -match "already exists") {
            Write-Host " Already exists." -ForegroundColor Yellow
        } else {
            Write-Host " Failed." -ForegroundColor Red
            Write-Host "Error: $errBody"
        }
    } else {
        Write-Host " Failed." -ForegroundColor Red
        Write-Host "Error: $($_.Exception.Message)"
    }
}

# Relationship: recipes -> categories (category_id -> categories.id)
$recipesToCategories = @{
    type = "pg_create_object_relationship"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "recipes"
        }
        name = "category"
        using = @{
            foreign_key_constraint_on = "category_id"
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Creating recipes -> categories relationship..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $recipesToCategories
    Write-Host " Success." -ForegroundColor Green
} catch {
    $stream = $_.Exception.Response.GetResponseStream()
    $reader = [System.IO.StreamReader]::new($stream)
    $errBody = $reader.ReadToEnd()
    if ($errBody -match "already exists") {
        Write-Host " Already exists." -ForegroundColor Yellow
    } else {
        Write-Host " Failed." -ForegroundColor Red
    }
}

# Relationship: users -> recipes (one-to-many)
$usersToRecipes = @{
    type = "pg_create_array_relationship"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "users"
        }
        name = "recipes"
        using = @{
            foreign_key_constraint_on = @{
                table = @{
                    schema = "public"
                    name = "recipes"
                }
                column = "user_id"
            }
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Creating users -> recipes relationship..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $usersToRecipes
    Write-Host " Success." -ForegroundColor Green
} catch {
    $stream = $_.Exception.Response.GetResponseStream()
    $reader = [System.IO.StreamReader]::new($stream)
    $errBody = $reader.ReadToEnd()
    if ($errBody -match "already exists") {
        Write-Host " Already exists." -ForegroundColor Yellow
    } else {
        Write-Host " Failed." -ForegroundColor Red
    }
}

# Relationship: recipes -> recipe_ingredients (one-to-many)
$recipesToIngredients = @{
    type = "pg_create_array_relationship"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "recipes"
        }
        name = "recipe_ingredients"
        using = @{
            foreign_key_constraint_on = @{
                table = @{
                    schema = "public"
                    name = "recipe_ingredients"
                }
                column = "recipe_id"
            }
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Creating recipes -> recipe_ingredients relationship..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $recipesToIngredients
    Write-Host " Success." -ForegroundColor Green
} catch {
    $stream = $_.Exception.Response.GetResponseStream()
    $reader = [System.IO.StreamReader]::new($stream)
    $errBody = $reader.ReadToEnd()
    if ($errBody -match "already exists") {
        Write-Host " Already exists." -ForegroundColor Yellow
    } else {
        Write-Host " Failed." -ForegroundColor Red
    }
}

Write-Host "`nRelationships configuration complete!" -ForegroundColor Cyan

