$baseUrl = "http://localhost:8081"
$email = "chef@example.com"
$password = "secret123"

try {
    # 1. Login
    Write-Host "1. Logging in..."
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body (@{ email=$email; password=$password } | ConvertTo-Json) -ContentType "application/json"
    $token = $loginResponse.token
    Write-Host "   Token received."

    # 2. Create
    Write-Host "2. Creating recipe..."
    $headers = @{ Authorization = "Bearer $token" }
    $createBody = @{
        category_id = 1
        title = "Manual Test Recipe"
        description = "Original Description"
        preparation_time = 10
        price = 5.00
        thumbnail_url = "http://img.com/1.jpg"
        ingredients = @(@{name="Ing1"; quantity="1"; unit="kg"})
        steps = @(@{instruction="Step1"; image_url=""})
    } | ConvertTo-Json -Depth 5
    $createResponse = Invoke-RestMethod -Uri "$baseUrl/recipes" -Method Post -Body $createBody -ContentType "application/json" -Headers $headers
    $id = $createResponse.id
    Write-Host "   Created Recipe ID: $id"

    # 3. Edit
    Write-Host "3. Editing recipe $id..."
    $editBody = @{
        category_id = 1
        title = "EDITED Manual Recipe"
        description = "Edited Description"
        preparation_time = 20
        price = 10.00
        thumbnail_url = "http://img.com/2.jpg"
        ingredients = @(@{name="Ing1"; quantity="2"; unit="kg"})
        steps = @(@{instruction="Step1"; image_url=""})
    } | ConvertTo-Json -Depth 5
    $editResponse = Invoke-RestMethod -Uri "$baseUrl/recipes/$id" -Method Put -Body $editBody -ContentType "application/json" -Headers $headers
    Write-Host "   Edit Response: $($editResponse.message)"

    # 4. Delete
    Write-Host "4. Deleting recipe $id..."
    $deleteResponse = Invoke-RestMethod -Uri "$baseUrl/recipes/$id" -Method Delete -Headers $headers
    Write-Host "   Delete Response: $($deleteResponse.message)"

} catch {
    Write-Error "Error: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        Write-Host "   Details: $($reader.ReadToEnd())"
    }
}
