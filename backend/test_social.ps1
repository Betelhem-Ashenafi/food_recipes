$baseUrl = "http://localhost:8081"
$timestamp = Get-Date -Format "yyyyMMddHHmmss"
$email = "social_$timestamp@example.com"
$password = "social123"

# 1. Signup/Login
Write-Host "1. Signup/Login..."
try {
    $signupBody = @{
        name = "Social User"
        email = $email
        password = $password
    } | ConvertTo-Json
    Invoke-RestMethod -Uri "$baseUrl/signup" -Method Post -Body $signupBody -ContentType "application/json" -ErrorAction Stop
    Write-Host "Signup successful."
} catch {
    Write-Host "Signup failed (maybe exists), proceeding to login..."
}

# Always login to get token
Write-Host "Logging in..."
$loginBody = @{
    email = $email
    password = $password
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $response.token
} catch {
    Write-Error "Login failed: $($_.Exception.Message)"
    exit
}

if (-not $token) {
    Write-Error "No token received!"
    exit
}
$headers = @{ Authorization = "Bearer $token" }
Write-Host "Logged in. Token: $token"

# 2. Create a Recipe to interact with
Write-Host "`n2. Creating a Recipe..."
$recipeBody = @{
    category_id = 1
    title = "Social Recipe"
    description = "A recipe to like and comment on"
    preparation_time = 15
    price = 10.00
    thumbnail_url = "http://example.com/social.jpg"
    ingredients = @(
        @{ name = "Test Ingredient"; quantity = "1"; unit = "pc" }
    )
    steps = @(
        @{ instruction = "Test Step"; image_url = "" }
    )
} | ConvertTo-Json -Depth 10

try {
    $recipeResponse = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Post -Body $recipeBody -Headers $headers -ContentType "application/json"
    $recipeId = $recipeResponse.id
    Write-Host "Recipe created with ID: $recipeId"
} catch {
    Write-Error "Create Recipe failed: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        Write-Host "Details: $($reader.ReadToEnd())"
    }
    exit
}

# 3. Like Recipe
Write-Host "`n3. Liking Recipe..."
Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/like" -Method Post -Headers $headers
Write-Host "Recipe liked."

# 4. Unlike Recipe
Write-Host "`n4. Unliking Recipe..."
Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/like" -Method Delete -Headers $headers
Write-Host "Recipe unliked."

# 5. Bookmark Recipe
Write-Host "`n5. Bookmarking Recipe..."
Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/bookmark" -Method Post -Headers $headers
Write-Host "Recipe bookmarked."

# 6. Unbookmark Recipe
Write-Host "`n6. Unbookmarking Recipe..."
Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/bookmark" -Method Delete -Headers $headers
Write-Host "Recipe unbookmarked."

# 7. Add Comment
Write-Host "`n7. Adding Comment..."
$commentBody = @{ content = "This looks delicious!" } | ConvertTo-Json
try {
    $commentResp = Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/comments" -Method Post -Body $commentBody -Headers $headers -ContentType "application/json"
    Write-Host "Response: $($commentResp | ConvertTo-Json -Depth 5)"
    Write-Host "Comment added. ID: $($commentResp.id)"
} catch {
    Write-Error "Add Comment failed: $($_.Exception.Message)"
}

Start-Sleep -Seconds 1

# 8. Get Comments
Write-Host "`n8. Fetching Comments..."
try {
    $comments = Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/comments" -Method Get
    Write-Host "Raw Comments: $($comments | ConvertTo-Json -Depth 5)"
    if ($comments.Count -gt 0) {
        Write-Host "Found $($comments.Count) comments:"
        $comments | ForEach-Object { Write-Host "- [$($_.user_name)]: $($_.content)" }
    } else {
        Write-Host "No comments found!"
    }
} catch {
    Write-Error "Get Comments failed: $($_.Exception.Message)"
}

# 9. Rate Recipe
Write-Host "`n9. Rating Recipe (5 stars)..."
$rateBody = @{ rating = 5 } | ConvertTo-Json
try {
    Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/rate" -Method Post -Body $rateBody -Headers $headers -ContentType "application/json"
    Write-Host "Recipe rated."
} catch {
    Write-Error "Rate Recipe failed: $($_.Exception.Message)"
}

# 10. Get Ratings
Write-Host "`n10. Fetching Ratings..."
try {
    $rating = Invoke-RestMethod -Uri "$baseUrl/recipes/$recipeId/rate" -Method Get
    Write-Host "Average Rating: $($rating.average_rating)"
    Write-Host "Count: $($rating.count)"
} catch {
    Write-Error "Get Ratings failed: $($_.Exception.Message)"
}


