$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8081"
$uploadUrl = "$baseUrl/upload"
$dummyFileName = "test_image.txt"
$dummyFileContent = "Hello World Upload " + (Get-Date).ToString("yyyyMMddHHmmss")
$dummyFilePath = Join-Path $PWD "test_image.txt"

Write-Host "Creating dummy file at $dummyFilePath..."
Set-Content -Path $dummyFilePath -Value $dummyFileContent

Write-Host "Uploading file using curl.exe..."
# Use curl.exe for reliable multipart/form-data
$responseJson = curl.exe -s -F "file=@$dummyFilePath" $uploadUrl

if (-not $responseJson) {
    Write-Error "Upload failed: No response from server."
}

Write-Host "Response: $responseJson"

try {
    $responseObj = $responseJson | ConvertFrom-Json
} catch {
    Write-Error "Failed to parse JSON response."
}

if ($responseObj.url) {
    $fileUrl = $responseObj.url
    Write-Host "File uploaded successfully. URL: $fileUrl"
    
    # Verify by downloading
    Write-Host "Verifying file content..."
    try {
        $downloadedContent = Invoke-RestMethod -Uri $fileUrl
        if ($downloadedContent -match "Hello World Upload") {
            Write-Host "Success: Content matches!" -ForegroundColor Green
        } else {
            Write-Error "Failure: Content mismatch. Got: '$downloadedContent'"
        }
    } catch {
        Write-Error "Failed to download file from $fileUrl. Error: $_"
    }

} else {
    Write-Error "Upload failed. Response did not contain 'url'."
}

