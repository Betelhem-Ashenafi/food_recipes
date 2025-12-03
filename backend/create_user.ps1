$baseUrl = "http://localhost:8081"
$signupUrl = "$baseUrl/signup"

$signupBody = @{
    email = "chef@example.com"
    password = "password123"
    name = "Chef John"
} | ConvertTo-Json

try {
    Invoke-RestMethod -Uri $signupUrl -Method Post -Body $signupBody -ContentType "application/json"
    Write-Host "User created successfully."
} catch {
    Write-Host "User might already exist or error: $_"
}
