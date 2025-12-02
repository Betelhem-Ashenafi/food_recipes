$baseUrl = "http://localhost:8081"
$email = "chef@example.com"
$password = "secret123"

# 1. Signup (Ignore error if already exists)
Write-Host "1. Signing up..."
$signupBody = @{
    name = "Chef Master"
    email = $email
    password = $password
} | ConvertTo-Json

try {
    $signupResponse = Invoke-RestMethod -Uri "$baseUrl/signup" -Method Post -Body $signupBody -ContentType "application/json"
    Write-Host "Signup successful: $($signupResponse.user.name)"
} catch {
    Write-Host "Signup failed or user exists: $($_.Exception.Message)"
}

# 2. Login
Write-Host "`n2. Logging in..."
$loginBody = @{
    email = $email
    password = $password
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    if (-not $token) {
        Write-Error "Login failed: No token received"
        exit
    }
    Write-Host "Login successful. Token received."
} catch {
    Write-Error "Login failed: $($_.Exception.Message)"
    exit
}

# 3. Create Recipe
Write-Host "`n3. Creating Recipe..."
$recipeBody = @{
    category_id = 1
    title = "Spaghetti Carbonara"
    description = "Classic Italian pasta dish"
    preparation_time = 20
    price = 15.50
    thumbnail_url = "http://example.com/carbonara.jpg"
    ingredients = @(
        @{ name = "Spaghetti"; quantity = "200"; unit = "g" },
        @{ name = "Eggs"; quantity = "2"; unit = "pcs" },
        @{ name = "Pancetta"; quantity = "100"; unit = "g" }
    )
    steps = @(
        @{ instruction = "Boil pasta"; image_url = "" },
        @{ instruction = "Fry pancetta"; image_url = "" },
        @{ instruction = "Mix eggs and cheese"; image_url = "" },
        @{ instruction = "Combine all"; image_url = "" }
    )
} | ConvertTo-Json -Depth 10

$headers = @{
    Authorization = "Bearer $token"
}

try {
    $createResponse = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Post -Body $recipeBody -ContentType "application/json" -Headers $headers
    Write-Host "Recipe created! ID: $($createResponse.id)"
} catch {
    Write-Error "Create Recipe failed: $($_.Exception.Message)"
    # Print detailed error if available
    if ($_.Exception.Response) {
        $stream = $_.Exception.Response.GetResponseStream()
        if ($stream) {
            $reader = New-Object System.IO.StreamReader($stream)
            Write-Host "Details: $($reader.ReadToEnd())"
        }
    }
}

# 4. List Recipes
Write-Host "`n4. Listing Recipes..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Get
    Write-Host "Found $($recipes.Count) recipes:"
    $recipes | ForEach-Object { Write-Host "- [$($_.id)] $($_.title) ($($_.price))" }
} catch {
    Write-Error "List Recipes failed: $($_.Exception.Message)"
}
