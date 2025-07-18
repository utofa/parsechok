package selenium

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	whatsappURL    = "https://web.whatsapp.com"
	qrCodeXPath    = "//*[@id='app']/div/div/div[2]/div[1]/div/div[2]/div/canvas"
	defaultTimeout = 30 * time.Second
	minPort        = 9515
	maxPort        = 9999
)

// WhatsAppClient handles WhatsApp Web automation
type WhatsAppClient struct {
	driver  selenium.WebDriver
	service *selenium.Service
}

// cleanOldSessions removes sessions older than 24 hours
func cleanOldSessions(baseDir string) error {
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		if time.Since(file.ModTime()) > 24*time.Second {
			path := filepath.Join(baseDir, file.Name())
			if err := os.RemoveAll(path); err != nil {
				log.Printf("Warning: failed to remove old session directory %s: %v", path, err)
			}
		}
	}
	return nil
}

// findFreePort finds an available port between minPort and maxPort
func findFreePort() (int, error) {
	for port := minPort; port <= maxPort; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		listener.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no free ports available between %d and %d", minPort, maxPort)
}

// waitForPort waits for a port to become available
func waitForPort(port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for port %d", port)
}

// NewWhatsAppClient creates a new WhatsApp automation client
func NewWhatsAppClient() (*WhatsAppClient, error) {
	log.Println("Initializing WhatsApp client...")

	// Find available port for ChromeDriver
	port, err := findFreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port: %v", err)
	}
	log.Printf("Using port %d for ChromeDriver\n", port)

	// Base directory for Chrome user data
	baseDir := filepath.Join(".", "chrome_data")

	// Clean old sessions
	if err := cleanOldSessions(baseDir); err != nil {
		log.Printf("Warning: failed to clean old sessions: %v", err)
	}

	// Create unique user data directory with timestamp and process ID
	timestamp := time.Now().Format("20060102_150405")
	sessionDir := fmt.Sprintf("session_%s_%d", timestamp, os.Getpid())
	userDataDir, err := filepath.Abs(filepath.Join(baseDir, sessionDir))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Remove directory if it exists
	if err := os.RemoveAll(userDataDir); err != nil {
		return nil, fmt.Errorf("failed to remove existing user data directory: %v", err)
	}

	// Create fresh directory
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %v", err)
	}

	// Use the existing ChromeDriver path
	chromeDriverPath := filepath.Join("driver", "chromedriver.exe")
	log.Printf("Using ChromeDriver at: %s\n", chromeDriverPath)

	// Configure Chrome options
	chromeOpts := chrome.Capabilities{
		Path: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
		Args: []string{
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-gpu",
			"--window-size=1920,1080",
			"--start-maximized",
			"--disable-notifications",
			"--disable-popup-blocking",
			"--disable-infobars",
			"--disable-extensions",
			"--disable-blink-features=AutomationControlled",
			fmt.Sprintf("--user-data-dir=%s", userDataDir),
			"--remote-debugging-port=0",
			"--no-first-run",
			"--no-default-browser-check",
			"--ignore-certificate-errors",
			"--test-type",
		},
		ExcludeSwitches: []string{
			"enable-automation",
			"enable-logging",
		},
		Prefs: map[string]interface{}{
			"profile.default_content_setting_values.notifications": 2,
			"credentials_enable_service":                           false,
			"profile.password_manager_enabled":                     false,
			"profile.default_content_settings.popups":              0,
		},
	}

	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
		selenium.Output(os.Stderr),
	}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// Add Chrome options to capabilities
	caps.AddChrome(chromeOpts)

	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start ChromeDriver: %v", err)
	}
	log.Println("ChromeDriver service started successfully")

	// Wait for ChromeDriver to be ready
	if err := waitForPort(port, 10*time.Second); err != nil {
		service.Stop()
		return nil, fmt.Errorf("ChromeDriver not ready: %v", err)
	}

	// Create WebDriver instance
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		service.Stop()
		return nil, fmt.Errorf("failed to create WebDriver: %v", err)
	}
	log.Println("WebDriver instance created successfully")

	// Set implicit wait timeout
	if err := driver.SetImplicitWaitTimeout(defaultTimeout); err != nil {
		driver.Quit()
		service.Stop()
		return nil, fmt.Errorf("failed to set implicit wait timeout: %v", err)
	}

	// Set page load timeout
	if err := driver.SetPageLoadTimeout(defaultTimeout); err != nil {
		driver.Quit()
		service.Stop()
		return nil, fmt.Errorf("failed to set page load timeout: %v", err)
	}

	return &WhatsAppClient{
		driver:  driver,
		service: service,
	}, nil
}

// waitForElement waits for an element to be present and visible
func (c *WhatsAppClient) waitForElement(by, value string, timeout time.Duration) (selenium.WebElement, error) {
	log.Printf("Waiting for element: %s=%s (timeout: %v)\n", by, value, timeout)
	deadline := time.Now().Add(timeout)
	var element selenium.WebElement
	var err error

	for time.Now().Before(deadline) {
		element, err = c.driver.FindElement(by, value)
		if err == nil {
			visible, err := element.IsDisplayed()
			if err == nil && visible {
				log.Printf("Element found and visible: %s=%s\n", by, value)
				return element, nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	log.Printf("Element not found or not visible: %s=%s\n", by, value)
	return nil, fmt.Errorf("element not found or not visible after %v", timeout)
}

// GetQRCode opens WhatsApp Web and returns the QR code element
func (c *WhatsAppClient) GetQRCode(sessionID string) (string, error) {
	log.Printf("Getting QR code for session %s...\n", sessionID)

	// Delete all cookies before starting
	if err := c.driver.DeleteAllCookies(); err != nil {
		log.Printf("Warning: failed to delete cookies: %v\n", err)
	}

	// Navigate to WhatsApp Web
	if err := c.driver.Get(whatsappURL); err != nil {
		return "", fmt.Errorf("failed to open WhatsApp Web: %v", err)
	}

	// Check if already authorized
	script := `
		return window.localStorage.getItem('WAToken1') !== null && 
			   window.localStorage.getItem('WAToken2') !== null;
	`
	result, err := c.driver.ExecuteScript(script, nil)
	if err == nil {
		if isAuthorized, ok := result.(bool); ok && isAuthorized {
			return "", fmt.Errorf("Already authorized")
		}
	}

	// Print current URL to verify we're on the right page
	currentURL, err := c.driver.CurrentURL()
	if err != nil {
		log.Printf("Error getting current URL: %v\n", err)
	} else {
		log.Printf("Current URL: %s\n", currentURL)
	}

	// Wait for QR code to appear with timeout
	log.Println("Waiting for QR code element...")
	qrElement, err := c.waitForElement(selenium.ByXPATH, qrCodeXPath, defaultTimeout)
	if err != nil {
		// Try alternative QR code selector
		log.Println("Trying alternative QR code selector...")
		qrElement, err = c.waitForElement(selenium.ByCSSSelector, "canvas", defaultTimeout)
		if err != nil {
			return "", fmt.Errorf("failed to find QR code element: %v", err)
		}
	}

	// Get QR code data URL
	log.Println("Getting QR code data...")
	dataURL, err := qrElement.GetAttribute("data-url")
	if err != nil {
		log.Println("Failed to get data-url attribute, trying src attribute...")
		// Try getting the QR code as an image source if data-url is not available
		dataURL, err = qrElement.GetAttribute("src")
		if err != nil {
			// Try to get canvas image data
			log.Println("Trying to get canvas image data...")
			script := `
				var canvas = arguments[0];
				return canvas.toDataURL('image/png');
			`
			result, err := c.driver.ExecuteScript(script, []interface{}{qrElement})
			if err != nil {
				return "", fmt.Errorf("failed to get QR code data: %v", err)
			}
			if dataURL, ok := result.(string); ok {
				return dataURL, nil
			}
			return "", fmt.Errorf("failed to convert QR code data to string")
		}
	}

	log.Printf("QR code data obtained (length: %d) for session %s\n", len(dataURL), sessionID)
	return dataURL, nil
}

// GetSessionData retrieves cookies and local storage data
func (c *WhatsAppClient) GetSessionData() ([]byte, error) {
	cookies, err := c.driver.GetCookies()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %v", err)
	}

	// Execute JavaScript to get localStorage
	localStorage, err := c.driver.ExecuteScript("return Object.entries(localStorage)", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get localStorage: %v", err)
	}

	sessionData := map[string]interface{}{
		"cookies":      cookies,
		"localStorage": localStorage,
	}

	return json.Marshal(sessionData)
}

// RestoreSession restores a previous session using cookies and localStorage
func (c *WhatsAppClient) RestoreSession(sessionData []byte) error {
	// First navigate to WhatsApp Web
	if err := c.driver.Get(whatsappURL); err != nil {
		return fmt.Errorf("failed to open WhatsApp Web: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(sessionData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal session data: %v", err)
	}

	// Restore cookies
	if cookies, ok := data["cookies"].([]selenium.Cookie); ok {
		for _, cookie := range cookies {
			if err := c.driver.AddCookie(&cookie); err != nil {
				return fmt.Errorf("failed to restore cookie: %v", err)
			}
		}
	}

	// Restore localStorage
	if localStorage, ok := data["localStorage"].([]interface{}); ok {
		script := "localStorage.clear();"
		for _, item := range localStorage {
			if pair, ok := item.([]interface{}); ok && len(pair) == 2 {
				key := pair[0].(string)
				value := pair[1].(string)
				script += fmt.Sprintf("localStorage.setItem('%s', '%s');", key, value)
			}
		}
		if _, err := c.driver.ExecuteScript(script, nil); err != nil {
			return fmt.Errorf("failed to restore localStorage: %v", err)
		}
	}

	// Refresh the page after restoring session data
	if err := c.driver.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh page: %v", err)
	}

	return nil
}

// SendMessage sends a message to a specific phone number
func (c *WhatsAppClient) SendMessage(phoneNumber, message string) error {
	// Open chat with phone number
	url := fmt.Sprintf("%s/send?phone=%s", whatsappURL, phoneNumber)
	if err := c.driver.Get(url); err != nil {
		return fmt.Errorf("failed to open chat: %v", err)
	}

	// Wait for message input to be ready
	input, err := c.waitForElement(selenium.ByCSSSelector, "div[contenteditable='true']", defaultTimeout)
	if err != nil {
		return fmt.Errorf("failed to find message input: %v", err)
	}

	if err := input.SendKeys(message); err != nil {
		return fmt.Errorf("failed to input message: %v", err)
	}

	// Send message
	if err := input.SendKeys(selenium.EnterKey); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// Close closes the WebDriver session and ChromeDriver service
func (c *WhatsAppClient) Close() error {
	if err := c.driver.Quit(); err != nil {
		return fmt.Errorf("failed to quit driver: %v", err)
	}
	c.service.Stop()
	return nil
}
