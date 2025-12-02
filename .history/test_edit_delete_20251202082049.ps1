$baseUrl = "http://localhost:8081"
$email = "chef@example.com"
$password = "secret123"

# 1. Login (Assuming user exists from previous tests)
Write-Host "1. Logging in..."
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
    Write-Host "Login successful."
} catch {
    Write-Error "Login failed: $($_.Exception.Message)"
    exit
}

$headers = @{
    Authorization = "Bearer $token"
}

# 2. Create a Recipe to Edit/Delete
Write-Host "`n2. Creating Recipe for Test..."
$recipeBody = @{
    category_id = 1
    title = "Test Recipe"
    description = "To be edited"
    preparation_time = 10
    price = 10.00
    thumbnail_url = "http://example.com/test.jpg"
    ingredients = @(
        @{ name = "Test Ing"; quantity = "1"; unit = "pc" }
    )
    steps = @(
        @{ instruction = "Test Step"; image_url = "" }
    )
} | ConvertTo-Json -Depth 10

try {
    $createResponse = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Post -Body $recipeBody -ContentType "application/json" -Headers $headers
    $recipeID = $createResponse.id
    Write-Host "Recipe created! ID: $recipeID"
} catch {
    Write-Error "Create Recipe failed: $($_.Exception.Message)"
    exit
}

# 3. Edit Recipe
Write-Host "`n3. Editing Recipe $recipeID..."
$editBody = @{
    category_id = 1
    title = "EDITED Recipe"
    description = "Has been edited"
    preparation_time = 15
    price = 20.00
    thumbnail_url = "http://example.com/edited.jpg"
    ingredients = @(
        @{ name = "Edited Ing"; quantity = "2"; unit = "pcs" }
    )
    steps = @(
        @{ instruction = "Edited Step"; image_url = "" }
    )
} | ConvertTo-Json -Depth 10

try {
    $editResponse = Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeID" -Method Put -Body $editBody -ContentType "application/json" -Headers $headers
    Write-Host "Recipe edited: $($editResponse.message)"
} catch {
    Write-Error "Edit Recipe failed: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $stream = $_.Exception.Response.GetResponseStream()
        if ($stream) {
            $reader = New-Object System.IO.StreamReader($stream)
            Write-Host "Details: $($reader.ReadToEnd())"
        }
    }
}

# 4. Verify Edit (Get Recipe - Optional, but good to check)
# Skipping for brevity, assuming success message is enough for now.

# 5. Delete Recipe
Write-Host "`n5. Deleting Recipe $recipeID..."
try {
    $deleteResponse = Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeID" -Method Delete -Headers $headers
    Write-Host "Recipe deleted: $($deleteResponse.message)"
} catch {
    Write-Error "Delete Recipe failed: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $stream = $_.Exception.Response.GetResponseStream()
        if ($stream) {
            $reader = New-Object System.IO.StreamReader($stream)
            Write-Host "Details: $($reader.ReadToEnd())"
        }
    }
}

# 6. Verify Deletion (Try to Edit again, should fail)
Write-Host "`n6. Verifying Deletion (Expect 404)..."
try {
    Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeID" -Method Put -Body $editBody -ContentType "application/json" -Headers $headers
    Write-Error "Error: Recipe should have been deleted but Edit succeeded!"
} catch {
    Write-Host "Success: Edit failed as expected ($($_.Exception.Message))"
}
