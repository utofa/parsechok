# Download and setup ChromeDriver
$chromeDriverVersion = "114.0.5735.90" # Update this version as needed
$downloadUrl = "https://chromedriver.storage.googleapis.com/$chromeDriverVersion/chromedriver_win32.zip"
$outputPath = ".\chromedriver.zip"
$extractPath = ".\bin"

# Create bin directory if it doesn't exist
New-Item -ItemType Directory -Force -Path $extractPath

# Download ChromeDriver
Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath

# Extract ChromeDriver
Expand-Archive -Path $outputPath -DestinationPath $extractPath -Force

# Add bin directory to PATH for this session
$env:PATH = "$env:PATH;$(Resolve-Path $extractPath)"

# Clean up zip file
Remove-Item $outputPath

Write-Host "ChromeDriver has been installed to $extractPath"
Write-Host "Please add this path to your system's PATH environment variable:"
Write-Host "$(Resolve-Path $extractPath)" 