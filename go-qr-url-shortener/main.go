package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

// Template parsing
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// In-memory storage for shortened URLs
var urlDatabase = make(map[string]string)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/generate-qr", generateQRPage)
	http.HandleFunc("/shorten-url", shortenURLPage)
	http.HandleFunc("/process-qr", processQR)
	http.HandleFunc("/process-url", processURL)
	http.HandleFunc("/redirect/", redirectToOriginalURL)

	// Serve static files (QR codes, CSS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Home Page
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// QR Code Page
func generateQRPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "qr.html", nil)
}

// URL Shortener Page
func shortenURLPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "short.html", nil)
}

// Process QR Code Generation
func processQR(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/generate-qr", http.StatusSeeOther)
		return
	}

	text := r.FormValue("text")
	qrPath := "static/qr_codes/qrcode.png"

	err := qrcode.WriteFile(text, qrcode.Medium, 256, qrPath)
	if err != nil {
		http.Error(w, "Error generating QR Code", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "qr.html", qrPath)
}

// Process URL Shortening
func processURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/shorten-url", http.StatusSeeOther)
		return
	}

	originalURL := r.FormValue("url")

	// Ensure URL starts with http:// or https://
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	// Generate short code
	shortCode := generateShortCode()
	urlDatabase[shortCode] = originalURL

	// Create shortened URL
	shortURL := fmt.Sprintf("http://localhost:8080/redirect/%s", shortCode)

	tmpl.ExecuteTemplate(w, "short.html", shortURL)
}

// Redirect from short URL to original URL
func redirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	// Get short code from URL
	shortCode := strings.TrimPrefix(r.URL.Path, "/redirect/")

	// Find original URL
	originalURL, exists := urlDatabase[shortCode]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// Generate a random short URL code
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	shortCode := make([]byte, 6)
	for i := range shortCode {
		shortCode[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortCode)
}
