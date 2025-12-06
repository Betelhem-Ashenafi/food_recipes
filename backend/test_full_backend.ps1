$baseUrl = "http://localhost:8081"
$email = "testuser_$(Get-Random)@example.com"
$password = "password123"

# --- 1. Signup ---
Write-Host "1. Testing Signup..."
$signupBody = @{
    name = "Test User"
    email = $email
    password = $password
} | ConvertTo-Json

try {
    $signupResponse = Invoke-WebRequest -Uri "$baseUrl/signup" -Method Post -Body $signupBody -ContentType "application/json"
    Write-Host "Signup Success: $($signupResponse.Content)"
} catch {
    Write-Host "Signup Failed: $_"
    exit
}

# --- 2. Login ---
Write-Host "`n2. Testing Login..."
$loginBody = @{
    email = $email
    password = $password
} | ConvertTo-Json

try {
    $loginResponse = Invoke-WebRequest -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $loginData = $loginResponse.Content | ConvertFrom-Json
    $token = $loginData.token
    Write-Host "Login Success. Token received."
} catch {
    Write-Host "Login Failed: $_"
    exit
}

$headers = @{
    Authorization = "Bearer $token"
}

# --- 3. Get Categories ---
Write-Host "`n3. Testing Get Categories..."
try {
    $catResponse = Invoke-WebRequest -Uri "$baseUrl/categories" -Method Get
    $categories = $catResponse.Content | ConvertFrom-Json
    Write-Host "Categories found: $($categories.Count)"
    if ($categories.Count -eq 0) {
        $categoryId = 1 
    } else {
        $categoryId = $categories[0].id
    }
} catch {
    Write-Host "Get Categories Failed: $_"
}

# --- 4. Create Recipe ---
Write-Host "`n4. Testing Create Recipe..."
$recipeBody = @{
    category_id = $categoryId
    title = "Full Test Recipe"
    description = "A delicious test recipe for full backend check"
    preparation_time = 45
    price = 25.00
    thumbnail_url = "http://example.com/image.jpg"
    ingredients = @(
        @{ name = "Flour"; quantity = "500"; unit = "g" },
        @{ name = "Sugar"; quantity = "200"; unit = "g" }
    )
    steps = @(
        @{ instruction = "Mix ingredients"; image_url = "" },
        @{ instruction = "Bake at 200C"; image_url = "" }
    )
} | ConvertTo-Json -Depth 10

try {
    $createResponse = Invoke-WebRequest -Uri "$baseUrl/recipes" -Method Post -Body $recipeBody -ContentType "application/json" -Headers $headers
    $createData = $createResponse.Content | ConvertFrom-Json
    $recipeId = $createData.id
    Write-Host "Create Recipe Success: ID $recipeId"
} catch {
    Write-Host "Create Recipe Failed: $_"
    exit
}

# --- 5. Edit Recipe ---
Write-Host "`n5. Testing Edit Recipe..."
$editBody = @{
    category_id = $categoryId
    title = "Updated Test Recipe"
    description = "Updated description"
    preparation_time = 50
    price = 30.00
    thumbnail_url = "http://example.com/image_updated.jpg"
    ingredients = @(
        @{ name = "Flour"; quantity = "600"; unit = "g" }
    )
    steps = @(
        @{ instruction = "Mix well"; image_url = "" }
    )
} | ConvertTo-Json -Depth 10

try {
    $editResponse = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId" -Method Put -Body $editBody -ContentType "application/json" -Headers $headers
    Write-Host "Edit Recipe Success: $($editResponse.Content)"
} catch {
    Write-Host "Edit Recipe Failed: $_"
}

# --- 6. Social Features (Like, Bookmark, Comment, Rate) ---
Write-Host "`n6. Testing Social Features..."

# Like
try {
    $likeResponse = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId/like" -Method Post -Headers $headers
    Write-Host "Like Recipe Success"
} catch {
    Write-Host "Like Recipe Failed: $_"
}

# Bookmark
try {
    $bookmarkResponse = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId/bookmark" -Method Post -Headers $headers
    Write-Host "Bookmark Recipe Success"
} catch {
    Write-Host "Bookmark Recipe Failed: $_"
}

# Comment
$commentBody = @{ content = "Great recipe!" } | ConvertTo-Json
try {
    $commentResponse = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId/comments" -Method Post -Body $commentBody -ContentType "application/json" -Headers $headers
    Write-Host "Comment Recipe Success"
} catch {
    Write-Host "Comment Recipe Failed: $_"
}

# Rate
$rateBody = @{ rating = 5 } | ConvertTo-Json
try {
    $rateResponse = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId/rate" -Method Post -Body $rateBody -ContentType "application/json" -Headers $headers
    Write-Host "Rate Recipe Success"
} catch {
    Write-Host "Rate Recipe Failed: $_"
}

# --- 7. Payment Initialization ---
Write-Host "`n7. Testing Payment Initialization..."
$paymentBody = @{
    amount = "100"
    email = $email
    first_name = "Test"
    last_name = "User"
    recipe_id = $recipeId
} | ConvertTo-Json

try {
    $paymentResponse = Invoke-WebRequest -Uri "$baseUrl/payment/initialize" -Method Post -Body $paymentBody -ContentType "application/json" -Headers $headers
    $paymentData = $paymentResponse.Content | ConvertFrom-Json
    $txRef = $paymentData.tx_ref
    Write-Host "Payment Init Success. TxRef: $txRef"
} catch {
    Write-Host "Payment Init Failed: $_"
}

# --- 8. Payment Verification (Mock) ---
# Note: This requires the mock Chapa server to be running on port 8082
Write-Host "`n8. Testing Payment Verification..."
try {
    $verifyResponse = Invoke-WebRequest -Uri "$baseUrl/payment/verify?tx_ref=$txRef" -Method Get -Headers $headers
    Write-Host "Payment Verify Success: $($verifyResponse.Content)"
} catch {
    Write-Host "Payment Verify Failed: $_"
}

# --- 9. Delete Recipe ---
Write-Host "`n9. Testing Delete Recipe..."
try {
    $deleteResponse = Invoke-WebRequest -Uri "$baseUrl/recipes/$recipeId" -Method Delete -Headers $headers
    Write-Host "Delete Recipe Success: $($deleteResponse.Content)"
} catch {
    Write-Host "Delete Recipe Failed: $_"
}

Write-Host "`n--- Full Backend Check Completed ---"
