# Simple Backend Test Script
$baseUrl = "http://localhost:8081"
$passed = 0
$failed = 0

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  SIMPLE BACKEND TEST" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Test 1: Signup
Write-Host "[1] Testing Signup..." -ForegroundColor Yellow
$email = "testuser_$(Get-Random)@example.com"
$password = "password123"
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/signup" -Method Post -Body (@{name="Test User"; email=$email; password=$password} | ConvertTo-Json) -ContentType "application/json" -ErrorAction Stop
    Write-Host "✅ Signup - PASSED" -ForegroundColor Green
    $passed++
} catch {
    Write-Host "❌ Signup - FAILED: $($_.Exception.Message)" -ForegroundColor Red
    $failed++
    exit
}

# Test 2: Login
Write-Host "`n[2] Testing Login..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/login" -Method Post -Body (@{email=$email; password=$password} | ConvertTo-Json) -ContentType "application/json" -ErrorAction Stop
    $token = ($response.Content | ConvertFrom-Json).token
    Write-Host "✅ Login - PASSED (Token received)" -ForegroundColor Green
    $passed++
} catch {
    Write-Host "❌ Login - FAILED: $($_.Exception.Message)" -ForegroundColor Red
    $failed++
    exit
}

# Test 3: Get Categories
Write-Host "`n[3] Testing Get Categories..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/categories" -Method Get -ErrorAction Stop
    $categories = $response.Content | ConvertFrom-Json
    Write-Host "✅ Get Categories - PASSED ($($categories.Count) categories found)" -ForegroundColor Green
    $categoryId = $categories[0].id
    $passed++
} catch {
    Write-Host "❌ Get Categories - FAILED: $($_.Exception.Message)" -ForegroundColor Red
    $failed++
    $categoryId = 1
}

# Test 4: Create Recipe
Write-Host "`n[4] Testing Create Recipe..." -ForegroundColor Yellow
$headers = @{Authorization = "Bearer $token"}
$recipeBody = @{
    category_id = $categoryId
    title = "Test Recipe"
    description = "Test description"
    preparation_time = 30
    price = 0
    thumbnail_url = "http://example.com/image.jpg"
    ingredients = @(@{name="Flour"; quantity="2"; unit="cups"})
    steps = @(@{instruction="Mix ingredients"})
}
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/recipes" -Method Post -Headers $headers -Body ($recipeBody | ConvertTo-Json -Depth 10) -ContentType "application/json" -ErrorAction Stop
    $recipeId = ($response.Content | ConvertFrom-Json).id
    Write-Host "✅ Create Recipe - PASSED (Recipe ID: $recipeId)" -ForegroundColor Green
    $passed++
} catch {
    Write-Host "❌ Create Recipe - FAILED: $($_.Exception.Message)" -ForegroundColor Red
    $failed++
}

# Test 5: Get Recipes
Write-Host "`n[5] Testing Get Recipes..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/recipes" -Method Get -ErrorAction Stop
    Write-Host "✅ Get Recipes - PASSED" -ForegroundColor Green
    $passed++
} catch {
    Write-Host "❌ Get Recipes - FAILED: $($_.Exception.Message)" -ForegroundColor Red
    $failed++
}

# Test 6: Get Ingredients
if ($recipeId) {
    Write-Host "`n[6] Testing Get Ingredients..." -ForegroundColor Yellow
    try {
        $response = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId/ingredients" -Method Get -ErrorAction Stop
        Write-Host "✅ Get Ingredients - PASSED" -ForegroundColor Green
        $passed++
    } catch {
        Write-Host "❌ Get Ingredients - FAILED: $($_.Exception.Message)" -ForegroundColor Red
        $failed++
    }

    # Test 7: Get Steps
    Write-Host "`n[7] Testing Get Steps..." -ForegroundColor Yellow
    try {
        $response = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId/steps" -Method Get -ErrorAction Stop
        Write-Host "✅ Get Steps - PASSED" -ForegroundColor Green
        $passed++
    } catch {
        Write-Host "❌ Get Steps - FAILED: $($_.Exception.Message)" -ForegroundColor Red
        $failed++
    }
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  TEST SUMMARY" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor $(if ($failed -eq 0) { "Green" } else { "Red" })
Write-Host "`n========================================`n" -ForegroundColor Cyan

if ($failed -eq 0) {
    Write-Host "🎉 ALL TESTS PASSED!" -ForegroundColor Green
} else {
    Write-Host "⚠️ Some tests failed" -ForegroundColor Yellow
}



