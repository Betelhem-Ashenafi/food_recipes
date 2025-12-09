# Setup Hasura Authentication Actions
$headers = @{
    "Content-Type" = "application/json"
    "X-Hasura-Admin-Secret" = "myhasurasecret"
}

Write-Host ""
Write-Host "Setting up Hasura Authentication Actions..." -ForegroundColor Cyan
Write-Host ""

# Step 1: Update custom types
Write-Host "[1] Updating custom types..." -ForegroundColor Yellow

$newTypes = @{
    type = "set_custom_types"
    args = @{
        scalars = @()
        input_objects = @(
            @{
                name = "LoginInput"
                fields = @(
                    @{ name = "email"; type = "String!" },
                    @{ name = "password"; type = "String!" }
                )
            },
            @{
                name = "SignupInput"
                fields = @(
                    @{ name = "name"; type = "String!" },
                    @{ name = "email"; type = "String!" },
                    @{ name = "password"; type = "String!" }
                )
            }
        )
        objects = @(
            @{
                name = "LoginOutput"
                fields = @(
                    @{ name = "token"; type = "String!" },
                    @{ name = "user_id"; type = "Int!" },
                    @{ name = "name"; type = "String!" },
                    @{ name = "email"; type = "String!" }
                )
            },
            @{
                name = "SignupOutput"
                fields = @(
                    @{ name = "id"; type = "Int!" },
                    @{ name = "name"; type = "String!" },
                    @{ name = "email"; type = "String!" }
                )
            }
        )
        enums = @()
    }
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $newTypes | Out-Null
Write-Host "  Custom types set" -ForegroundColor Green
Write-Host ""

# Step 2: Create login action
Write-Host "[2] Creating login action..." -ForegroundColor Yellow

$loginAction = @{
    type = "create_action"
    args = @{
        name = "login"
        definition = @{
            kind = "synchronous"
            arguments = @(
                @{ name = "arg"; type = "LoginInput!" }
            )
            output_type = "LoginOutput!"
            handler = "http://host.docker.internal:8081/hasura/login"
            timeout = 30
        }
    }
} | ConvertTo-Json -Depth 10

try {
    Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $loginAction | Out-Null
    Write-Host "  Login action created" -ForegroundColor Green
} catch {
    Write-Host "  Login action already exists" -ForegroundColor Yellow
}
Write-Host ""

# Step 3: Create signup action
Write-Host "[3] Creating signup action..." -ForegroundColor Yellow

$signupAction = @{
    type = "create_action"
    args = @{
        name = "signup"
        definition = @{
            kind = "synchronous"
            arguments = @(
                @{ name = "arg"; type = "SignupInput!" }
            )
            output_type = "SignupOutput!"
            handler = "http://host.docker.internal:8081/hasura/signup"
            timeout = 30
        }
    }
} | ConvertTo-Json -Depth 10

try {
    Invoke-RestMethod -Uri "http://localhost:8080/v1/metadata" -Method Post -Headers $headers -Body $signupAction | Out-Null
    Write-Host "  Signup action created" -ForegroundColor Green
} catch {
    Write-Host "  Signup action already exists" -ForegroundColor Yellow
}
Write-Host ""

Write-Host "======================================" -ForegroundColor Green
Write-Host "  HASURA AUTH ACTIONS READY!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host ""
