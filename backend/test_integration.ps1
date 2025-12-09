# Integration Test Script for Food Recipes App
# Tests all features: Social, Profile, Edit/Delete, Payment

Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                                                               ║" -ForegroundColor Cyan
Write-Host "║   🧪 FULL SYSTEM INTEGRATION TEST 🧪                         ║" -ForegroundColor Cyan -BackgroundColor DarkCyan
Write-Host "║                                                               ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════════════════════╝`n" -ForegroundColor Cyan

$baseUrl = "http://localhost:8081"
$testResults = @()

# Helper function to test endpoint
function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Url,
        [hashtable]$Headers = @{},
        [string]$Body = $null,
        [int]$ExpectedStatus = 200
    )
    
    Write-Host "`n[TEST] $Name" -ForegroundColor Yellow
    Write-Host "  $Method $Url" -ForegroundColor Gray
    
    try {
        $params = @{
            Uri = $Url
            Method = $Method
            Headers = $Headers
            UseBasicParsing = $true
            ErrorAction = "Stop"
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = "application/json"
        }
        
        $response = Invoke-WebRequest @params
        $statusCode = $response.StatusCode
        
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "  ✅ PASS - Status: $statusCode" -ForegroundColor Green
            $script:testResults += @{Name=$Name; Status="PASS"; StatusCode=$statusCode}
            return $response
        } else {
            Write-Host "  ⚠️  UNEXPECTED - Status: $statusCode (Expected: $ExpectedStatus)" -ForegroundColor Yellow
            $script:testResults += @{Name=$Name; Status="WARN"; StatusCode=$statusCode}
            return $response
        }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "  ✅ PASS - Status: $statusCode" -ForegroundColor Green
            $script:testResults += @{Name=$Name; Status="PASS"; StatusCode=$statusCode}
        } else {
            Write-Host "  ❌ FAIL - Status: $statusCode (Expected: $ExpectedStatus)" -ForegroundColor Red
            Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
            $script:testResults += @{Name=$Name; Status="FAIL"; StatusCode=$statusCode; Error=$_.Exception.Message}
        }
        return $null
    }
}

# Step 1: Test Signup
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "STEP 1: USER SIGNUP" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

$testEmail = "test$(Get-Date -Format 'yyyyMMddHHmmss')@test.com"
$testPassword = "test123456"
$signupData = @{
    name = "Test User"
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

$signupResponse = Test-Endpoint -Name "User Signup" -Method "POST" -Url "$baseUrl/signup" -Body $signupData -ExpectedStatus 200

if ($signupResponse) {
    $signupResult = $signupResponse.Content | ConvertFrom-Json
    Write-Host "  User created: $($signupResult.user.email)" -ForegroundColor Green
}

# Step 2: Test Login
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "STEP 2: USER LOGIN" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

$loginData = @{
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

$loginResponse = Test-Endpoint -Name "User Login" -Method "POST" -Url "$baseUrl/login" -Body $loginData -ExpectedStatus 200

if ($loginResponse) {
    $loginResult = $loginResponse.Content | ConvertFrom-Json
    $token = $loginResult.token
    $userId = $loginResult.user.id
    Write-Host "  ✅ Token obtained: $($token.Substring(0, 30))..." -ForegroundColor Green
    Write-Host "  User ID: $userId" -ForegroundColor Green
} else {
    Write-Host "  ❌ Cannot continue without token" -ForegroundColor Red
    exit 1
}

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Step 3: Test Social Features
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "STEP 3: SOCIAL FEATURES (Like/Comment/Rate/Bookmark)" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

# Get a recipe ID (assuming recipe ID 1 exists)
$recipeId = 1

# Test Like
Test-Endpoint -Name "Like Recipe" -Method "POST" -Url "$baseUrl/recipes/$recipeId/like" -Headers $headers -ExpectedStatus 200

# Test Check Like
Test-Endpoint -Name "Check Like" -Method "GET" -Url "$baseUrl/recipes/$recipeId/like/check" -Headers $headers -ExpectedStatus 200

# Test Bookmark
Test-Endpoint -Name "Bookmark Recipe" -Method "POST" -Url "$baseUrl/recipes/$recipeId/bookmark" -Headers $headers -ExpectedStatus 200

# Test Check Bookmark
Test-Endpoint -Name "Check Bookmark" -Method "GET" -Url "$baseUrl/recipes/$recipeId/bookmark/check" -Headers $headers -ExpectedStatus 200

# Test Comment
$commentData = @{
    content = "This is a test comment from integration test!"
} | ConvertTo-Json
Test-Endpoint -Name "Post Comment" -Method "POST" -Url "$baseUrl/recipes/$recipeId/comments" -Headers $headers -Body $commentData -ExpectedStatus 201

# Test Get Comments
Test-Endpoint -Name "Get Comments" -Method "GET" -Url "$baseUrl/recipes/$recipeId/comments" -ExpectedStatus 200

# Test Rating
$ratingData = @{
    rating = 5
} | ConvertTo-Json
Test-Endpoint -Name "Rate Recipe" -Method "POST" -Url "$baseUrl/recipes/$recipeId/rate" -Headers $headers -Body $ratingData -ExpectedStatus 200

# Test Get Rating
Test-Endpoint -Name "Get Rating" -Method "GET" -Url "$baseUrl/recipes/$recipeId/rate" -ExpectedStatus 200

# Step 4: Test Profile Endpoints
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "STEP 4: USER PROFILE ENDPOINTS" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

Test-Endpoint -Name "Get User Bookmarks" -Method "GET" -Url "$baseUrl/users/$userId/bookmarks" -Headers $headers -ExpectedStatus 200

Test-Endpoint -Name "Get User Purchases" -Method "GET" -Url "$baseUrl/users/$userId/purchases" -Headers $headers -ExpectedStatus 200

# Step 5: Test Recipe Edit/Delete
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "STEP 5: RECIPE EDIT/DELETE (Ownership Check)" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

# First, create a recipe to test edit/delete
Write-Host "`n[TEST] Create Recipe for Edit/Delete Test" -ForegroundColor Yellow
$createRecipeData = @{
    title = "Test Recipe for Edit/Delete"
    description = "This recipe will be edited and deleted"
    category_id = 1
    preparation_time = 30
    price = 0
    ingredients = @(
        @{name = "Test Ingredient"; quantity = "1"; unit = "cup"}
    )
    steps = @(
        @{instruction = "Test step 1"; image_url = ""}
    )
} | ConvertTo-Json

$createResponse = Test-Endpoint -Name "Create Recipe" -Method "POST" -Url "$baseUrl/recipes" -Headers $headers -Body $createRecipeData -ExpectedStatus 201

if ($createResponse) {
    $createResult = $createResponse.Content | ConvertFrom-Json
    $testRecipeId = $createResult.id
    Write-Host "  ✅ Recipe created with ID: $testRecipeId" -ForegroundColor Green
    
    # Test Edit Recipe
    $editRecipeData = @{
        title = "Updated Test Recipe"
        description = "This recipe has been updated"
        category_id = 1
        preparation_time = 45
        price = 0
        ingredients = @(
            @{name = "Updated Ingredient"; quantity = "2"; unit = "cups"}
        )
        steps = @(
            @{instruction = "Updated step 1"; image_url = ""}
        )
    } | ConvertTo-Json
    
    Test-Endpoint -Name "Edit Recipe" -Method "PUT" -Url "$baseUrl/recipes/$testRecipeId" -Headers $headers -Body $editRecipeData -ExpectedStatus 200
    
    # Test Delete Recipe
    Test-Endpoint -Name "Delete Recipe" -Method "DELETE" -Url "$baseUrl/recipes/$testRecipeId" -Headers $headers -ExpectedStatus 200
}

# Step 6: Test Payment Endpoints
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "STEP 6: PAYMENT SYSTEM (Chapa)" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

$paymentData = @{
    amount = "100"
    email = $testEmail
    first_name = "Test"
    last_name = "User"
    recipe_id = $recipeId
} | ConvertTo-Json

$paymentResponse = Test-Endpoint -Name "Initialize Payment" -Method "POST" -Url "$baseUrl/payment/initialize" -Headers $headers -Body $paymentData -ExpectedStatus 200

if ($paymentResponse) {
    $paymentResult = $paymentResponse.Content | ConvertFrom-Json
    Write-Host "  ✅ Payment initialized" -ForegroundColor Green
    Write-Host "  Checkout URL: $($paymentResult.checkout_url)" -ForegroundColor Gray
    Write-Host "  TX Ref: $($paymentResult.tx_ref)" -ForegroundColor Gray
}

# Step 7: Test Purchase Check
Test-Endpoint -Name "Check Purchase" -Method "GET" -Url "$baseUrl/recipes/$recipeId/purchase/check" -Headers $headers -ExpectedStatus 200

# Summary
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "TEST SUMMARY" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

$passed = ($testResults | Where-Object { $_.Status -eq "PASS" }).Count
$failed = ($testResults | Where-Object { $_.Status -eq "FAIL" }).Count
$warned = ($testResults | Where-Object { $_.Status -eq "WARN" }).Count
$total = $testResults.Count

Write-Host "`nTotal Tests: $total" -ForegroundColor White
Write-Host "✅ Passed: $passed" -ForegroundColor Green
Write-Host "⚠️  Warnings: $warned" -ForegroundColor Yellow
Write-Host "❌ Failed: $failed" -ForegroundColor Red

if ($failed -eq 0) {
    Write-Host "`n🎉 ALL TESTS PASSED!" -ForegroundColor Green -BackgroundColor DarkGreen
} else {
    Write-Host "`n⚠️  Some tests failed. Check the output above." -ForegroundColor Yellow
}

Write-Host "`n"

