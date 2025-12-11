$ErrorActionPreference = "Continue"
$baseUrl = "http://localhost:8081"

Write-Host "`n=== CHAPA PAYMENT INTEGRATION TEST ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Login
Write-Host "STEP 1: User Login" -ForegroundColor Yellow
$loginBody = @{email="test@example.com"; password="password123"} | ConvertTo-Json
try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    $userId = $loginResponse.user.id
    Write-Host "[OK] Login successful - User ID: $userId" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Login failed: $_" -ForegroundColor Red
    exit 1
}

# Step 2: Initialize Payment
Write-Host "`nSTEP 2: Initialize Payment" -ForegroundColor Yellow
$headers = @{"Authorization"="Bearer $token"; "Content-Type"="application/json"}
$paymentBody = @{
    amount = "100"
    email = "testuser@gmail.com"
    first_name = "Test"
    last_name = "User"
    recipe_id = 1
} | ConvertTo-Json

try {
    $paymentResponse = Invoke-RestMethod -Uri "$baseUrl/payment/initialize" -Method Post -Headers $headers -Body $paymentBody
    $checkoutUrl = $paymentResponse.checkout_url
    $txRef = $paymentResponse.tx_ref
    Write-Host "[OK] Payment initialized!" -ForegroundColor Green
    Write-Host "     Checkout URL: $checkoutUrl" -ForegroundColor Cyan
    Write-Host "     Transaction Ref: $txRef" -ForegroundColor Cyan
    Write-Host "`n[INFO] Open the checkout URL in browser to complete test payment" -ForegroundColor Yellow
} catch {
    Write-Host "[ERROR] Payment initialization failed!" -ForegroundColor Red
    Write-Host "        Error: $_" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "        Response: $errorBody" -ForegroundColor Red
    }
}

Write-Host "`n=== TEST COMPLETE ===" -ForegroundColor Cyan
