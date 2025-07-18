# WhatsApp Parser

A Golang-based WhatsApp parser using Selenium WebDriver for automating WhatsApp Web interactions.

## Features
- QR code generation for WhatsApp Web authentication
- Session management (save/restore)
- Message sending functionality
- Clean architecture implementation

## Requirements
- Go 1.20+
- Chrome/Chromium browser
- ChromeDriver
- Selenium WebDriver

## Setup
1. Install dependencies:
```bash
go mod download
```

2. Make sure ChromeDriver is installed and in your PATH

3. Run the application:
```bash
go run cmd/app/main.go
```

## Project Structure
```
.
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   └── delivery/
│       └── http/
├── pkg/
│   └── selenium/
└── config/
``` 