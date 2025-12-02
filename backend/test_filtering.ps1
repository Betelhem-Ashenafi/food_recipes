$baseUrl = "http://localhost:8081"

# 1. Login to get token
$loginBody = @{
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    Write-Host "Login successful. Token: $token"
} catch {
    Write-Error "Login failed: $_"
    exit
}

$headers = @{
    Authorization = "Bearer $token"
}

# Helper function to print recipes
function Print-Recipes($recipes, $label) {
    Write-Host "--- $label ---"
    if ($recipes.Count -eq 0) {
        Write-Host "No recipes found."
    } else {
        foreach ($r in $recipes) {
            Write-Host "ID: $($r.id), Title: $($r.title), Time: $($r.preparation_time)"
        }
    }
    Write-Host ""
}

# 2. Test Filter by Title
Write-Host "Testing Filter by Title (e.g., 'Pasta')..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?title=Pasta" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes with 'Pasta' in title"
} catch {
    Write-Error "Filter by title failed: $_"
}

# 3. Test Filter by Time
Write-Host "Testing Filter by Time (<= 30 mins)..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?time=30" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes <= 30 mins"
} catch {
    Write-Error "Filter by time failed: $_"
}

# 4. Test Filter by Ingredient
Write-Host "Testing Filter by Ingredient (e.g., 'Tomato')..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?ingredient=Tomato" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes with 'Tomato'"
} catch {
    Write-Error "Filter by ingredient failed: $_"
}

# 5. Test Filter by Creator
Write-Host "Testing Filter by Creator (e.g., 'Test User')..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?creator=Test" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes by 'Test'"
} catch {
    Write-Error "Filter by creator failed: $_"
}

# 6. Test Combined Filter
Write-Host "Testing Combined Filter (Title='Pasta' AND Time<=45)..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?title=Pasta&time=45" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes 'Pasta' <= 45 mins"
} catch {
    Write-Error "Combined filter failed: $_"
}
