# Comprehensive Backend Testing Script
# Tests all backend functionality according to requirements

$baseUrl = "http://localhost:8081"
$testResults = @()

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Url,
        [hashtable]$Headers = @{},
        [object]$Body = $null,
        [int]$ExpectedStatus = 200
    )
    
    Write-Host "`n[TEST] $Name" -ForegroundColor Cyan
    Write-Host "  Method: $Method | URL: $Url" -ForegroundColor Gray
    
    try {
        $params = @{
            Uri = $Url
            Method = $Method
            Headers = $Headers
            ErrorAction = "Stop"
        }
        
        if ($Body) {
            $params.Body = ($Body | ConvertTo-Json -Depth 10)
            $params.ContentType = "application/json"
        }
        
        $response = Invoke-WebRequest @params
        $statusCode = $response.StatusCode
        
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "  ✓ PASSED (Status: $statusCode)" -ForegroundColor Green
            $script:testResults += @{
                Test = $Name
                Status = "PASSED"
                StatusCode = $statusCode
            }
            return $response
        } else {
            Write-Host "  ✗ FAILED (Expected: $ExpectedStatus, Got: $statusCode)" -ForegroundColor Red
            $script:testResults += @{
                Test = $Name
                Status = "FAILED"
                StatusCode = $statusCode
                Expected = $ExpectedStatus
            }
            return $null
        }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "  ✗ FAILED (Status: $statusCode)" -ForegroundColor Red
        Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Yellow
        $script:testResults += @{
            Test = $Name
            Status = "FAILED"
            StatusCode = $statusCode
            Error = $_.Exception.Message
        }
        return $null
    }
}

Write-Host "========================================" -ForegroundColor Magenta
Write-Host "  BACKEND COMPREHENSIVE TEST SUITE" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta

# ============================================
# 1. AUTHENTICATION TESTS
# ============================================
Write-Host "`n[1] AUTHENTICATION TESTS" -ForegroundColor Yellow

# Test Signup
$testEmail = "testuser_$(Get-Random)@example.com"
$testPassword = "password123"
$signupBody = @{
    name = "Test User"
    email = $testEmail
    password = $testPassword
}
$signupResponse = Test-Endpoint -Name "User Signup" -Method "POST" -Url "$baseUrl/signup" -Body $signupBody -ExpectedStatus 200

# Test Login
$loginBody = @{
    email = $testEmail
    password = $testPassword
}
$loginResponse = Test-Endpoint -Name "User Login" -Method "POST" -Url "$baseUrl/login" -Body $loginBody -ExpectedStatus 200

if ($loginResponse) {
    $loginData = $loginResponse.Content | ConvertFrom-Json
    $token = $loginData.token
    Write-Host "  Token received: $($token.Substring(0, [Math]::Min(20, $token.Length)))..." -ForegroundColor Gray
} else {
    Write-Host "  ✗ Cannot continue without authentication token" -ForegroundColor Red
    exit 1
}

$authHeaders = @{
    Authorization = "Bearer $token"
}

# ============================================
# 2. CATEGORIES TESTS
# ============================================
Write-Host "`n[2] CATEGORIES TESTS" -ForegroundColor Yellow

$categoriesResponse = Test-Endpoint -Name "Get Categories" -Method "GET" -Url "$baseUrl/categories"
if ($categoriesResponse) {
    $categories = $categoriesResponse.Content | ConvertFrom-Json
    $categoryId = $categories[0].id
    Write-Host "  Found $($categories.Count) categories, using category ID: $categoryId" -ForegroundColor Gray
} else {
    $categoryId = 1
    Write-Host "  Using fallback category ID: $categoryId" -ForegroundColor Yellow
}

# ============================================
# 3. FILE UPLOAD TESTS
# ============================================
Write-Host "`n[3] FILE UPLOAD TESTS" -ForegroundColor Yellow

# Create a test file
$testFileContent = "Test image content"
$testFilePath = "$env:TEMP\test_image.txt"
$testFileContent | Out-File -FilePath $testFilePath -Encoding utf8

try {
    $formData = @{
        file = Get-Item $testFilePath
    }
    $uploadResponse = Invoke-WebRequest -Uri "$baseUrl/upload" -Method POST -Form $formData -Headers $authHeaders -ErrorAction Stop
    if ($uploadResponse.StatusCode -eq 200) {
        $uploadData = $uploadResponse.Content | ConvertFrom-Json
        $uploadedUrl = $uploadData.url
        Write-Host "  ✓ File Upload PASSED" -ForegroundColor Green
        Write-Host "  Uploaded URL: $uploadedUrl" -ForegroundColor Gray
        $script:testResults += @{Test = "File Upload"; Status = "PASSED"; StatusCode = 200}
    }
} catch {
    Write-Host "  ✗ File Upload FAILED" -ForegroundColor Red
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Yellow
    $script:testResults += @{Test = "File Upload"; Status = "FAILED"; Error = $_.Exception.Message}
    $uploadedUrl = "http://localhost:8081/uploads/test.jpg"
}

Remove-Item $testFilePath -ErrorAction SilentlyContinue

# ============================================
# 4. RECIPE CRUD TESTS
# ============================================
Write-Host "`n[4] RECIPE CRUD TESTS" -ForegroundColor Yellow

# Create Recipe
$recipeBody = @{
    category_id = $categoryId
    title = "Test Recipe $(Get-Date -Format 'HHmmss')"
    description = "This is a test recipe description"
    preparation_time = 30
    price = 5.99
    thumbnail_url = $uploadedUrl
    ingredients = @(
        @{ name = "Flour"; quantity = "2"; unit = "cups" },
        @{ name = "Sugar"; quantity = "1"; unit = "cup" }
    )
    steps = @(
        @{ instruction = "Mix flour and sugar"; image_url = "" },
        @{ instruction = "Bake at 350F for 30 minutes"; image_url = "" }
    )
}

$createResponse = Test-Endpoint -Name "Create Recipe" -Method "POST" -Url "$baseUrl/recipes" -Headers $authHeaders -Body $recipeBody -ExpectedStatus 200

if ($createResponse) {
    $createData = $createResponse.Content | ConvertFrom-Json
    $recipeId = $createData.id
    Write-Host "  Recipe created with ID: $recipeId" -ForegroundColor Gray
} else {
    Write-Host "  ✗ Cannot continue without recipe ID" -ForegroundColor Red
    exit 1
}

# Get Recipes (with filters)
Test-Endpoint -Name "Get All Recipes" -Method "GET" -Url "$baseUrl/recipes"
Test-Endpoint -Name "Get Recipes by Title" -Method "GET" -Url "$baseUrl/recipes?title=Test"
Test-Endpoint -Name "Get Recipes by Time" -Method "GET" -Url "$baseUrl/recipes?time=60"
Test-Endpoint -Name "Get Recipes by Ingredient" -Method "GET" -Url "$baseUrl/recipes?ingredient=Flour"
Test-Endpoint -Name "Get Recipes by Creator" -Method "GET" -Url "$baseUrl/recipes?creator=Test"

# Edit Recipe
$editBody = @{
    category_id = $categoryId
    title = "Updated Test Recipe"
    description = "Updated description"
    preparation_time = 45
    price = 7.99
    thumbnail_url = $uploadedUrl
    ingredients = @(
        @{ name = "Flour"; quantity = "3"; unit = "cups" }
    )
    steps = @(
        @{ instruction = "Updated step"; image_url = "" }
    )
}
Test-Endpoint -Name "Edit Recipe" -Method "PUT" -Url "$baseUrl/recipes/$recipeId" -Headers $authHeaders -Body $editBody -ExpectedStatus 200

# ============================================
# 5. SOCIAL FEATURES TESTS
# ============================================
Write-Host "`n[5] SOCIAL FEATURES TESTS" -ForegroundColor Yellow

# Like Recipe
Test-Endpoint -Name "Like Recipe" -Method "POST" -Url "$baseUrl/recipes/$recipeId/like" -Headers $authHeaders -ExpectedStatus 200

# Bookmark Recipe
Test-Endpoint -Name "Bookmark Recipe" -Method "POST" -Url "$baseUrl/recipes/$recipeId/bookmark" -Headers $authHeaders -ExpectedStatus 200

# Comment on Recipe
$commentBody = @{ content = "This is a test comment" }
Test-Endpoint -Name "Comment on Recipe" -Method "POST" -Url "$baseUrl/recipes/$recipeId/comments" -Headers $authHeaders -Body $commentBody -ExpectedStatus 201

# Get Comments
Test-Endpoint -Name "Get Comments" -Method "GET" -Url "$baseUrl/recipes/$recipeId/comments"

# Rate Recipe
$rateBody = @{ rating = 5 }
Test-Endpoint -Name "Rate Recipe" -Method "POST" -Url "$baseUrl/recipes/$recipeId/rate" -Headers $authHeaders -Body $rateBody -ExpectedStatus 200

# Get Rating
Test-Endpoint -Name "Get Recipe Rating" -Method "GET" -Url "$baseUrl/recipes/$recipeId/rate"

# Unlike Recipe
Test-Endpoint -Name "Unlike Recipe" -Method "DELETE" -Url "$baseUrl/recipes/$recipeId/like" -Headers $authHeaders -ExpectedStatus 200

# Unbookmark Recipe
Test-Endpoint -Name "Unbookmark Recipe" -Method "DELETE" -Url "$baseUrl/recipes/$recipeId/bookmark" -Headers $authHeaders -ExpectedStatus 200

# ============================================
# 6. PAYMENT TESTS
# ============================================
Write-Host "`n[6] PAYMENT TESTS" -ForegroundColor Yellow

$paymentBody = @{
    amount = "10.00"
    email = $testEmail
    first_name = "Test"
    last_name = "User"
    recipe_id = $recipeId
}
$paymentResponse = Test-Endpoint -Name "Initialize Payment" -Method "POST" -Url "$baseUrl/payment/initialize" -Headers $authHeaders -Body $paymentBody -ExpectedStatus 200

if ($paymentResponse) {
    $paymentData = $paymentResponse.Content | ConvertFrom-Json
    $txRef = $paymentData.tx_ref
    Write-Host "  Payment initialized, tx_ref: $txRef" -ForegroundColor Gray
    Write-Host "  Note: Payment verification requires actual Chapa transaction" -ForegroundColor Yellow
}

# ============================================
# 7. DELETE RECIPE TEST
# ============================================
Write-Host "`n[7] DELETE RECIPE TEST" -ForegroundColor Yellow

Test-Endpoint -Name "Delete Recipe" -Method "DELETE" -Url "$baseUrl/recipes/$recipeId" -Headers $authHeaders -ExpectedStatus 200

# ============================================
# 8. HASURA INTEGRATION TESTS
# ============================================
Write-Host "`n[8] HASURA INTEGRATION TESTS" -ForegroundColor Yellow

# Test Hasura Login Action (Hasura format)
$hasuraLoginBody = @{
    action = @{
        name = "login"
    }
    input = @{
        arg = @{
            email = $testEmail
            password = $testPassword
        }
    }
}
Test-Endpoint -Name "Hasura Login Action" -Method "POST" -Url "$baseUrl/hasura/login" -Body $hasuraLoginBody -ExpectedStatus 200

# Test Event Trigger (would be called by Hasura, but we can test the endpoint)
$eventBody = @{
    event = @{
        op = "INSERT"
        data = @{
            new = @{
                id = 1
                title = "Test Recipe"
            }
        }
    }
    table = @{
        schema = "public"
        name = "recipes"
    }
}
Test-Endpoint -Name "Hasura Event Trigger" -Method "POST" -Url "$baseUrl/events/new-recipe" -Body $eventBody -ExpectedStatus 200

# Test Hasura Upload Action
Write-Host "`n[TEST] Hasura Upload Action" -ForegroundColor Cyan
Write-Host "  Method: POST | URL: $baseUrl/hasura/upload" -ForegroundColor Gray
try {
    $testFileContent = "Test image content for Hasura"
    $testFilePath = "$env:TEMP\test_hasura_upload.txt"
    $testFileContent | Out-File -FilePath $testFilePath -Encoding utf8
    
    $formData = @{
        file = Get-Item $testFilePath
    }
    $hasuraUploadResponse = Invoke-WebRequest -Uri "$baseUrl/hasura/upload" -Method POST -Form $formData -ErrorAction Stop
    if ($hasuraUploadResponse.StatusCode -eq 200) {
        $hasuraUploadData = $hasuraUploadResponse.Content | ConvertFrom-Json
        Write-Host "  ✓ Hasura Upload Action PASSED" -ForegroundColor Green
        Write-Host "  Uploaded URL: $($hasuraUploadData.url)" -ForegroundColor Gray
        $script:testResults += @{Test = "Hasura Upload Action"; Status = "PASSED"; StatusCode = 200}
    }
    Remove-Item $testFilePath -ErrorAction SilentlyContinue
} catch {
    Write-Host "  ✗ Hasura Upload Action FAILED" -ForegroundColor Red
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Yellow
    $script:testResults += @{Test = "Hasura Upload Action"; Status = "FAILED"; Error = $_.Exception.Message}
}

# ============================================
# SUMMARY
# ============================================
Write-Host "`n========================================" -ForegroundColor Magenta
Write-Host "  TEST SUMMARY" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta

$passed = ($testResults | Where-Object { $_.Status -eq "PASSED" }).Count
$failed = ($testResults | Where-Object { $_.Status -eq "FAILED" }).Count
$total = $testResults.Count

Write-Host "`nTotal Tests: $total" -ForegroundColor White
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor $(if ($failed -eq 0) { "Green" } else { "Red" })

if ($failed -gt 0) {
    Write-Host "`nFailed Tests:" -ForegroundColor Red
    $testResults | Where-Object { $_.Status -eq "FAILED" } | ForEach-Object {
        Write-Host "  - $($_.Test)" -ForegroundColor Yellow
        if ($_.Error) {
            Write-Host "    Error: $($_.Error)" -ForegroundColor Gray
        }
    }
}

Write-Host "`n========================================" -ForegroundColor Magenta

