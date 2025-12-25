# Convert COPY statements to INSERT for SQL editor compatibility
$inputFile = "database_backup_clean.sql"
$outputFile = "database_backup_for_editor.sql"

Write-Host "Converting COPY to INSERT statements..." -ForegroundColor Yellow

$content = Get-Content $inputFile -Raw
$lines = $content -split "`n"

$output = @()
$inCopy = $false
$copyTable = ""
$copyColumns = ""

foreach ($line in $lines) {
    if ($line -match "^COPY\s+(\S+)\s+\((.+)\)\s+FROM\s+stdin;") {
        $inCopy = $true
        $copyTable = $matches[1]
        $copyColumns = $matches[2]
        continue
    }
    
    if ($line -match "^\\\.$") {
        $inCopy = $false
        continue
    }
    
    if ($inCopy -and $line.Trim() -ne "") {
        # Convert COPY data line to INSERT
        $values = $line -split "`t"
        $values = $values | ForEach-Object { 
            if ($_ -eq '\N' -or $_ -eq '') { 
                'NULL' 
            } else { 
                "'" + ($_ -replace "'", "''") + "'" 
            } 
        }
        $valuesStr = $values -join ", "
        $output += "INSERT INTO $copyTable ($copyColumns) VALUES ($valuesStr);"
    } else {
        if (-not ($line -match "^\\\.$")) {
            $output += $line
        }
    }
}

$output -join "`n" | Out-File $outputFile -Encoding UTF8 -NoNewline

Write-Host "✅ Created: $outputFile" -ForegroundColor Green
Write-Host "Size: $((Get-Item $outputFile).Length) characters" -ForegroundColor Cyan

