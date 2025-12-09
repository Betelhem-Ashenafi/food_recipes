# Configure Hasura Actions (Login & Signup)
$headers = @{
    "Content-Type" = "application/json"
    "X-Hasura-Admin-Secret" = "myhasurasecret"
}

$hasuraUrl = "http://localhost:8080"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  CONFIGURING HASURA ACTIONS" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Step 1: Define custom types
Write-Host "[1] Creating custom types..." -ForegroundColor Yellow

$typesPayload = @{
    type = "set_custom_types"
    args = @{
        scalars = @()
        input_objects = @(
            @{
                name = "LoginInput"
                fields = @(
                    @{
                        name = "email"
                        type = "String!"
                    },
                    @{
                        name = "password"
                        type = "String!"
                    }
                )
            },
            @{
                name = "SignupInput"
                fields = @(
                    @{
                        name = "name"
                        type = "String!"
                    },
                    @{
                        name = "email"
                        type = "String!"
                    },
                    @{
                        name = "password"
                        type = "String!"
                    }
                )
            }
        )
        objects = @(
            @{
                name = "LoginOutput"
                fields = @(
                    @{
                        name = "token"
                        type = "String!"
                    },
                    @{
                        name = "user_id"
                        type = "Int!"
                    },
                    @{
                        name = "name"
                        type = "String!"
                    },
                    @{
                        name = "email"
                        type = "String!"
                    }
                )
            },
            @{
                name = "SignupOutput"
                fields = @(
                    @{
                        name = "id"
                        type = "Int!"
                    },
                    @{
                        name = "name"
                        type = "String!"
                    },
                    @{
                        name = "email"
                        type = "String!"
                    }
                )
            }
        )
        enums = @()
    }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $typesPayload
    Write-Host "  ✅ Custom types created successfully`n" -ForegroundColor Green
} catch {
    Write-Host "  ⚠️  Custom types may already exist (continuing...)`n" -ForegroundColor Yellow
}

# Step 2: Create Login Action
Write-Host "[2] Creating login action..." -ForegroundColor Yellow

$loginActionPayload = @{
    type = "create_action"
    args = @{
        name = "login"
        definition = @{
            kind = "synchronous"
            arguments = @(
                @{
                    name = "arg"
                    type = "LoginInput!"
                }
            )
            output_type = "LoginOutput!"
            handler = "http://host.docker.internal:8081/hasura/login"
            timeout = 30
        }
        comment = "User login action"
    }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $loginActionPayload
    Write-Host "  ✅ Login action created successfully`n" -ForegroundColor Green
} catch {
    $errMsg = $_.Exception.Message
    if ($errMsg -match "already exists") {
        Write-Host "  ⚠️  Login action already exists, updating..." -ForegroundColor Yellow
        
        # Drop existing action
        $dropPayload = @{
            type = "drop_action"
            args = @{
                name = "login"
            }
        } | ConvertTo-Json -Depth 10
        
        try {
            Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $dropPayload | Out-Null
            # Recreate
            Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $loginActionPayload | Out-Null
            Write-Host "  ✅ Login action updated successfully`n" -ForegroundColor Green
        } catch {
            Write-Host "  ❌ Failed to update login action" -ForegroundColor Red
            Write-Host "  Error: $_`n" -ForegroundColor Red
        }
    } else {
        Write-Host "  ❌ Failed to create login action" -ForegroundColor Red
        Write-Host "  Error: $_`n" -ForegroundColor Red
    }
}

# Step 3: Create Signup Action
Write-Host "[3] Creating signup action..." -ForegroundColor Yellow

$signupActionPayload = @{
    type = "create_action"
    args = @{
        name = "signup"
        definition = @{
            kind = "synchronous"
            arguments = @(
                @{
                    name = "arg"
                    type = "SignupInput!"
                }
            )
            output_type = "SignupOutput!"
            handler = "http://host.docker.internal:8081/hasura/signup"
            timeout = 30
        }
        comment = "User signup action"
    }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $signupActionPayload
    Write-Host "  ✅ Signup action created successfully`n" -ForegroundColor Green
} catch {
    $errMsg = $_.Exception.Message
    if ($errMsg -match "already exists") {
        Write-Host "  ⚠️  Signup action already exists, updating..." -ForegroundColor Yellow
        
        # Drop existing action
        $dropPayload = @{
            type = "drop_action"
            args = @{
                name = "signup"
            }
        } | ConvertTo-Json -Depth 10
        
        try {
            Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $dropPayload | Out-Null
            # Recreate
            Invoke-RestMethod -Uri "$hasuraUrl/v1/metadata" -Method Post -Headers $headers -Body $signupActionPayload | Out-Null
            Write-Host "  ✅ Signup action updated successfully`n" -ForegroundColor Green
        } catch {
            Write-Host "  ❌ Failed to update signup action" -ForegroundColor Red
            Write-Host "  Error: $_`n" -ForegroundColor Red
        }
    } else {
        Write-Host "  ❌ Failed to create signup action" -ForegroundColor Red
        Write-Host "  Error: $_`n" -ForegroundColor Red
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  ✅ HASURA ACTIONS CONFIGURED!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Cyan
