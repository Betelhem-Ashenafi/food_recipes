$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8081"
$loginUrl = "$baseUrl/login"
$initPaymentUrl = "$baseUrl/payment/initialize"

# 1. Login to get token
$loginBody = @{
    email = "chapa_test_user@example.com"
    password = "password123"
} | ConvertTo-Json

Write-Host "Logging in..."
try {
    $loginResponse = Invoke-RestMethod -Uri $loginUrl -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    Write-Host "Login successful. Token received." -ForegroundColor Green
} catch {
    Write-Error "Login failed: $_"
}

# 2. Initialize Payment
$paymentBody = @{
    amount = "100"
    email = "abebe@gmail.com"
    first_name = "John"
    last_name = "Doe"
    recipe_id = 1
} | ConvertTo-Json

$headers = @{
    "Authorization" = "Bearer $token"
}

Write-Host "Initializing payment..."
try {
    $paymentResponse = Invoke-RestMethod -Uri $initPaymentUrl -Method Post -Headers $headers -Body $paymentBody -ContentType "application/json"
    Write-Host "Payment Initialized!" -ForegroundColor Green
    Write-Host "Checkout URL: $($paymentResponse.checkout_url)"
    Write-Host "Tx Ref: $($paymentResponse.tx_ref)"
} catch {
    Write-Host "Payment initialization failed (Expected if CHAPA_SECRET_KEY is invalid)." -ForegroundColor Yellow
    Write-Host "Error: $_"
    # Print the error details if available
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        Write-Host "Server Response: $($reader.ReadToEnd())"
    }
}
