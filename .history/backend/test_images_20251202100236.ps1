
# Test Script for Recipe Images
$ErrorActionPreference = "Stop"

$baseUrl = "http://localhost:8081"
$email = "testimage@example.com"
$password = "password123"

Write-Host "1. Signup/Login..."
try {
    $signupBody = @{
        name = "Image Tester"
        email = $email
        password = $password
    } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$baseUrl/signup" -Method Post -Body $signupBody -ContentType "application/json" -ErrorAction SilentlyContinue
} catch {
    # If signup fails, try login
    $loginBody = @{
        email = $email
        password = $password
    } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
}

$token = $response.token
$headers = @{
    Authorization = "Bearer $token"
}
Write-Host "Logged in. Token obtained."

Write-Host "`n2. Creating a Recipe..."
$recipeBody = @{
    category_id = 1
    title = "Recipe with Images"
    description = "Testing multiple images"
    preparation_time = 30
    price = 15.50
    thumbnail_url = "http://example.com/thumb.jpg"
    ingredients = @(
        @{ name = "Flour"; quantity = "2"; unit = "cups" }
    )
    steps = @(
        @{ instruction = "Mix"; image_url = "" }
    )
} | ConvertTo-Json -Depth 5

$recipeResponse = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Post -Body $recipeBody -Headers $headers -ContentType "application/json"
$recipeId = $recipeResponse.id
Write-Host "Recipe created with ID: $recipeId"

Write-Host "`n3. Uploading Images..."
$imagesBody = @{
    images = @(
        "http://example.com/img1.jpg",
        "http://example.com/img2.jpg",
        "http://example.com/img3.jpg"
    )
} | ConvertTo-Json

Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/images" -Method Post -Body $imagesBody -Headers $headers -ContentType "application/json"
Write-Host "Images uploaded."

Write-Host "`n4. Verifying Images in Database..."
# We can't easily check DB directly from here without sql tools, but we can check if the endpoint didn't fail.
# Ideally, we would have a GET endpoint for recipe details that includes images, but for now we assume success if 201 Created.

Write-Host "`n5. Setting Featured Image..."
# We need to know the Image IDs. Since we just inserted them, let's assume they are sequential or we'd need a GET endpoint.
# For this test, we'll try to feature the first image we just added. 
# Since we don't have the IDs returned in the upload response (it just says "success"), 
# we might need to query the DB or add a GET endpoint to be sure.
# However, let's try to fetch the recipe list and see if images are included (if the model was updated).

$recipes = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Get -Headers $headers
$myRecipe = $recipes | Where-Object { $_.id -eq $recipeId }

if ($myRecipe.images) {
    Write-Host "Images found in recipe response:"
    $myRecipe.images | Format-Table id, url, is_featured
    
    $imageIdToFeature = $myRecipe.images[1].id # Let's feature the second one
    Write-Host "Featuring Image ID: $imageIdToFeature"

    Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/images/$imageIdToFeature/feature" -Method Post -Headers $headers
    Write-Host "Featured image set."

    Write-Host "`n6. Verifying Featured Status..."
    $recipes = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Get -Headers $headers
    $myRecipe = $recipes | Where-Object { $_.id -eq $recipeId }
    $featuredImg = $myRecipe.images | Where-Object { $_.id -eq $imageIdToFeature }
    
    if ($featuredImg.is_featured) {
        Write-Host "SUCCESS: Image $imageIdToFeature is marked as featured!"
    } else {
        Write-Error "FAILURE: Image $imageIdToFeature is NOT marked as featured."
    }
} else {
    Write-Warning "No images returned in GET /recipes response. Did you update the GetRecipesHandler to join/fetch images?"
}
