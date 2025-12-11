$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8081"
$hasuraUrl = "http://localhost:8080/v1/graphql"

Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                                                               ║" -ForegroundColor Cyan
Write-Host "║     🧪 CHAPA PAYMENT INTEGRATION TEST (SANDBOX MODE) 🧪      ║" -ForegroundColor Cyan
Write-Host "║                                                               ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════════════════════╝`n" -ForegroundColor Cyan

# Step 1: Login to get JWT token
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "STEP 1: USER LOGIN" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

$loginBody = @{
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    $userId = $loginResponse.user.id
    Write-Host "✅ Login successful" -ForegroundColor Green
    Write-Host "   User ID: $userId" -ForegroundColor Gray
    Write-Host "   Token: $($token.Substring(0, 20))..." -ForegroundColor Gray
} catch {
    Write-Host "❌ Login failed: $_" -ForegroundColor Red
    Write-Host "   Creating test user first..." -ForegroundColor Yellow
    
    # Try to create user
    $signupBody = @{
        name = "Test User"
        email = "test@example.com"
        password = "password123"
    } | ConvertTo-Json
    
    try {
        $signupResponse = Invoke-RestMethod -Uri "$baseUrl/signup" -Method Post -Body $signupBody -ContentType "application/json"
        $token = $signupResponse.token
        $userId = $signupResponse.user.id
        Write-Host "✅ User created and logged in" -ForegroundColor Green
    } catch {
        Write-Host "❌ Failed to create user: $_" -ForegroundColor Red
        exit 1
    }
}

# Step 2: Find or create a premium recipe
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "STEP 2: FIND/CREATE PREMIUM RECIPE" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Query Hasura for recipes with price > 0
$graphqlQuery = @"
{
  recipes(where: {price: {_gt: 0}}, limit: 1) {
    id
    title
    price
    owner_id
  }
}
"@

$graphqlBody = @{
    query = $graphqlQuery
} | ConvertTo-Json

try {
    $hasuraHeaders = @{
        "Content-Type" = "application/json"
        "X-Hasura-Admin-Secret" = "myhasurasecret"
    }
    
    $recipeResponse = Invoke-RestMethod -Uri $hasuraUrl -Method Post -Headers $hasuraHeaders -Body $graphqlBody -ErrorAction Stop
    $recipe = $recipeResponse.data.recipes[0]
    
    if ($recipe) {
        Write-Host "✅ Found premium recipe:" -ForegroundColor Green
        Write-Host "   Recipe ID: $($recipe.id)" -ForegroundColor Gray
        Write-Host "   Title: $($recipe.title)" -ForegroundColor Gray
        Write-Host "   Price: $($recipe.price) ETB" -ForegroundColor Gray
        $recipeId = $recipe.id
        $recipePrice = $recipe.price
    } else {
        Write-Host "⚠️  No premium recipe found. Creating one..." -ForegroundColor Yellow
        
        # Create a premium recipe
        $createRecipeBody = @{
            title = "Premium Test Recipe - Chapa Integration"
            description = "This is a test premium recipe for Chapa payment testing"
            category_id = 1
            prep_time = 30
            price = 100.00
        } | ConvertTo-Json
        
        try {
            $createResponse = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Post -Headers $headers -Body $createRecipeBody
            $recipeId = $createResponse.recipe_id
            $recipePrice = 100.00
            Write-Host "✅ Premium recipe created (ID: $recipeId)" -ForegroundColor Green
        } catch {
            Write-Host "❌ Failed to create recipe: $_" -ForegroundColor Red
            exit 1
        }
    }
} catch {
    Write-Host "❌ Failed to query recipes: $_" -ForegroundColor Red
    exit 1
}

# Step 3: Initialize Payment
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "STEP 3: INITIALIZE CHAPA PAYMENT" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

$paymentBody = @{
    amount = $recipePrice.ToString()
    email = "test@example.com"
    first_name = "Test"
    last_name = "User"
    recipe_id = $recipeId
} | ConvertTo-Json

Write-Host "Payment Request:" -ForegroundColor Cyan
Write-Host "   Amount: $recipePrice ETB" -ForegroundColor Gray
Write-Host "   Recipe ID: $recipeId" -ForegroundColor Gray
Write-Host "   Email: test@example.com" -ForegroundColor Gray

try {
    $paymentResponse = Invoke-RestMethod -Uri "$baseUrl/payment/initialize" -Method Post -Headers $headers -Body $paymentBody
    $checkoutUrl = $paymentResponse.checkout_url
    $txRef = $paymentResponse.tx_ref
    
    Write-Host "✅ Payment initialized successfully!" -ForegroundColor Green
    Write-Host "   Checkout URL: $checkoutUrl" -ForegroundColor Cyan
    Write-Host "   Transaction Reference: $txRef" -ForegroundColor Cyan
    
    # Step 4: Verify Payment (Simulate - In real test, user would complete payment on Chapa)
    Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
    Write-Host "STEP 4: VERIFY PAYMENT (SIMULATED)" -ForegroundColor Yellow
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
    
    Write-Host "⚠️  Note: In sandbox mode, you need to:" -ForegroundColor Yellow
    Write-Host "   1. Open the checkout URL in a browser" -ForegroundColor Gray
    Write-Host "   2. Complete the payment using Chapa test credentials" -ForegroundColor Gray
    Write-Host "   3. After payment, Chapa will redirect to /payment/success?tx_ref=$txRef" -ForegroundColor Gray
    Write-Host "   4. The frontend will call /payment/verify to record the purchase" -ForegroundColor Gray
    Write-Host ""
    Write-Host "   For automated testing, we'll verify the payment endpoint structure..." -ForegroundColor Cyan
    
    # Test verification endpoint structure (won't work without actual payment)
    try {
        $verifyResponse = Invoke-RestMethod -Uri "$baseUrl/payment/verify?tx_ref=$txRef" -Method Get -Headers $headers -ErrorAction SilentlyContinue
        Write-Host "✅ Verification endpoint is accessible" -ForegroundColor Green
    } catch {
        if ($_.Exception.Response.StatusCode -eq 400) {
            Write-Host "⚠️  Verification endpoint requires actual payment completion" -ForegroundColor Yellow
            Write-Host "   This is expected - payment must be completed on Chapa first" -ForegroundColor Gray
        } else {
            Write-Host "❌ Verification endpoint error: $_" -ForegroundColor Red
        }
    }
    
} catch {
    Write-Host "❌ Payment initialization failed!" -ForegroundColor Red
    $errorDetails = $_.Exception.Response
    if ($errorDetails) {
        $reader = [System.IO.StreamReader]::new($errorDetails.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "   Error Response: $errorBody" -ForegroundColor Red
    } else {
        Write-Host "   Error: $_" -ForegroundColor Red
    }
    
    # Check if it's a Chapa API error
    if ($_ -match "CHASECK|chapa|invalid") {
        Write-Host "`n⚠️  CHAPA API ERROR DETECTED" -ForegroundColor Yellow
        Write-Host "   This might be due to:" -ForegroundColor Yellow
        Write-Host "   - Invalid test secret key" -ForegroundColor Gray
        Write-Host "   - Chapa API endpoint issue" -ForegroundColor Gray
        Write-Host "   - Network connectivity" -ForegroundColor Gray
        Write-Host ""
        Write-Host "   Current test key: CHASECK_TEST-LR1J8py5LqhoMlJdtVT5piJnWtJ66RZk" -ForegroundColor Cyan
        Write-Host "   (This is a placeholder - replace with your actual Chapa test key)" -ForegroundColor Yellow
    }
    exit 1
}

# Step 5: Check Database for Purchase Record (if payment was completed)
Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
Write-Host "STEP 5: DATABASE VERIFICATION" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

Write-Host "To verify purchase in database, run this SQL query:" -ForegroundColor Cyan
Write-Host ""
Write-Host "   SELECT * FROM purchases WHERE user_id = $userId AND recipe_id = $recipeId;" -ForegroundColor White
Write-Host ""
Write-Host "Or check by transaction reference:" -ForegroundColor Cyan
Write-Host "   SELECT * FROM purchases WHERE chapa_tx_ref = '$txRef';" -ForegroundColor White

# Step 6: Summary
Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║                                                               ║" -ForegroundColor Green
Write-Host "║              ✅ TEST SUMMARY ✅                                ║" -ForegroundColor Green
Write-Host "║                                                               ║" -ForegroundColor Green
Write-Host "╚═══════════════════════════════════════════════════════════════╝`n" -ForegroundColor Green

Write-Host "✅ Payment initialization: SUCCESS" -ForegroundColor Green
Write-Host "✅ Checkout URL generated: $checkoutUrl" -ForegroundColor Green
Write-Host "✅ Transaction reference: $txRef" -ForegroundColor Green
Write-Host ""
Write-Host "📋 NEXT STEPS:" -ForegroundColor Yellow
Write-Host "   1. Open checkout URL in browser: $checkoutUrl" -ForegroundColor White
Write-Host "   2. Use Chapa test credentials to complete payment" -ForegroundColor White
Write-Host "   3. After payment, verify purchase in database" -ForegroundColor White
Write-Host "   4. Check that recipe is unlocked for the user" -ForegroundColor White
Write-Host ""
Write-Host "🔑 CHAPA TEST CREDENTIALS:" -ForegroundColor Cyan
Write-Host "   - Use test card numbers from Chapa documentation" -ForegroundColor Gray
Write-Host "   - Test mode doesn't charge real money" -ForegroundColor Gray
Write-Host "   - Payment will be simulated and verified" -ForegroundColor Gray
Write-Host ""

