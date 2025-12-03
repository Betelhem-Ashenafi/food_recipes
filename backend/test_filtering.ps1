$baseUrl = "http://localhost:8081"

# 1. Signup/Login to get token
$email = "testfilter@example.com"
$password = "password123"

try {
    $signupBody = @{
        name = "Filter Tester"
        email = $email
        password = $password
    } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$baseUrl/signup" -Method Post -Body $signupBody -ContentType "application/json" -ErrorAction SilentlyContinue
    Write-Host "Signup successful."
} catch {
    Write-Host "Signup failed (user might exist), trying login..."
}

try {
    $loginBody = @{
        email = $email
        password = $password
    } | ConvertTo-Json
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
    if ($null -eq $recipes) {
        Write-Host "Result is null."
        return
    }
    
    # Check if it's an array or single object
    $count = 0
    if ($recipes -is [System.Array] -or $recipes -is [System.Collections.IEnumerable]) {
        $count = $recipes.Count
    } else {
        $count = 1
        $recipes = @($recipes) # Wrap in array
    }

    if ($count -eq 0) {
        Write-Host "No recipes found."
    } else {
        Write-Host "Found $count recipes."
        # Print raw for debugging if needed, or just the fields
        # Write-Host ($recipes | ConvertTo-Json -Depth 2) 
        foreach ($r in $recipes) {
            Write-Host "ID: $($r.id), Title: $($r.title), Time: $($r.preparation_time)"
        }
    }
    Write-Host ""
}

# 2. Test Filter by Title
Write-Host "Testing Filter by Title (e.g., 'Spaghetti')..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?title=Spaghetti" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes with 'Spaghetti' in title"
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
Write-Host "Testing Filter by Ingredient (e.g., 'Egg')..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?ingredient=Egg" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes with 'Egg'"
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
Write-Host "Testing Combined Filter (Title='Spaghetti' AND Time<=25)..."
try {
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes?title=Spaghetti&time=25" -Method Get -Headers $headers
    Print-Recipes $recipes "Recipes 'Spaghetti' <= 25 mins"
} catch {
    Write-Error "Combined filter failed: $_"
}
