$baseUrl = "http://localhost:8081"
$email = "testuser_$(Get-Random)@example.com"
$password = "password123"

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
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        Write-Host "Error Details: $($reader.ReadToEnd())"
    }
    exit
}

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
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        Write-Host "Error Details: $($reader.ReadToEnd())"
    }
    exit
}

Write-Host "`n3. Testing Get Categories..."
try {
    $catResponse = Invoke-WebRequest -Uri "$baseUrl/categories" -Method Get
    $categories = $catResponse.Content | ConvertFrom-Json
    Write-Host "Categories found: $($categories.Count)"
    if ($categories.Count -eq 0) {
        Write-Host "Warning: No categories found. Recipe creation might fail."
        $categoryId = 1 # Fallback
    } else {
        $categoryId = $categories[0].id
    }
} catch {
    Write-Host "Get Categories Failed: $_"
}

Write-Host "`n4. Testing Create Recipe..."
$recipeBody = @{
    category_id = $categoryId
    title = "Test Recipe"
    description = "A delicious test recipe"
    preparation_time = 30
    price = 15.50
    thumbnail_url = "http://example.com/image.jpg"
    ingredients = @(
        @{ name = "Flour"; quantity = "500"; unit = "g" },
        @{ name = "Water"; quantity = "200"; unit = "ml" }
    )
    steps = @(
        @{ instruction = "Mix ingredients"; image_url = "" },
        @{ instruction = "Bake"; image_url = "" }
    )
} | ConvertTo-Json -Depth 10

$headers = @{
    Authorization = "Bearer $token"
}

try {
    $createResponse = Invoke-WebRequest -Uri "$baseUrl/recipes" -Method Post -Body $recipeBody -ContentType "application/json" -Headers $headers
    Write-Host "Create Recipe Success: $($createResponse.Content)"
} catch {
    Write-Host "Create Recipe Failed: $_"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        Write-Host "Error Details: $($reader.ReadToEnd())"
    }
}

Write-Host "`n5. Testing Get Recipes..."
try {
    $getRecipesResponse = Invoke-WebRequest -Uri "$baseUrl/recipes" -Method Get
    Write-Host "Get Recipes Success. Count: $(($getRecipesResponse.Content | ConvertFrom-Json).Count)"
} catch {
    Write-Host "Get Recipes Failed: $_"
}
