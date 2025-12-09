# Configure Hasura Permissions
$headers = @{
    "Content-Type" = "application/json"
    "X-Hasura-Admin-Secret" = "myhasurasecret"
}

Write-Host "Configuring Hasura Permissions..." -ForegroundColor Cyan

# Public Role Permissions for Recipes (Read Only)
$publicRecipesPermissions = @{
    type = "pg_create_select_permission"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "recipes"
        }
        role = "public"
        permission = @{
            columns = "*"
            filter = @{}
            allow_aggregations = $true
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Setting public read permissions for recipes..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $publicRecipesPermissions
    Write-Host " Success." -ForegroundColor Green
} catch {
    $stream = $_.Exception.Response.GetResponseStream()
    $reader = [System.IO.StreamReader]::new($stream)
    $errBody = $reader.ReadToEnd()
    if ($errBody -match "already exists") {
        Write-Host " Already exists." -ForegroundColor Yellow
    } else {
        Write-Host " Failed." -ForegroundColor Red
        Write-Host "Error: $errBody"
    }
}

# Public Role Permissions for Categories (Read Only)
$publicCategoriesPermissions = @{
    type = "pg_create_select_permission"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "categories"
        }
        role = "public"
        permission = @{
            columns = "*"
            filter = @{}
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Setting public read permissions for categories..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $publicCategoriesPermissions
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

# Authenticated Role Permissions for Recipes (Create, Read, Update own, Delete own)
$authRecipesInsert = @{
    type = "pg_create_insert_permission"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "recipes"
        }
        role = "user"
        permission = @{
            columns = "*"
            check = @{
                user_id = @{
                    _eq = "X-Hasura-User-Id"
                }
            }
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Setting authenticated insert permissions for recipes..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $authRecipesInsert
    Write-Host " Success." -ForegroundColor Green
} catch {
    $stream = $_.Exception.Response.GetResponseStream()
    $reader = [System.IO.StreamReader]::new($stream)
    $errBody = $reader.ReadToEnd()
    if ($errBody -match "already exists") {
        Write-Host " Already exists." -ForegroundColor Yellow
    } else {
        Write-Host " Failed." -ForegroundColor Red
        Write-Host "Error: $errBody"
    }
}

$authRecipesUpdate = @{
    type = "pg_create_update_permission"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "recipes"
        }
        role = "user"
        permission = @{
            columns = "*"
            filter = @{
                user_id = @{
                    _eq = "X-Hasura-User-Id"
                }
            }
            check = @{
                user_id = @{
                    _eq = "X-Hasura-User-Id"
                }
            }
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Setting authenticated update permissions for recipes (own only)..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $authRecipesUpdate
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

$authRecipesDelete = @{
    type = "pg_create_delete_permission"
    args = @{
        source = "default"
        table = @{
            schema = "public"
            name = "recipes"
        }
        role = "user"
        permission = @{
            filter = @{
                user_id = @{
                    _eq = "X-Hasura-User-Id"
                }
            }
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Setting authenticated delete permissions for recipes (own only)..." -NoNewline
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $authRecipesDelete
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

Write-Host "`nPermissions configuration complete!" -ForegroundColor Cyan
Write-Host "Note: You may need to configure JWT claims mapping in Hasura console." -ForegroundColor Yellow

