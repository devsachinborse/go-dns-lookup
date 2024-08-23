package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Declare allRecords 
var allRecords []DNSRecord

// DNSLookupHandler handles DNS lookups and stores results in allRecords
func DNSLookupHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	domainInput := r.FormValue("domains")
	domains := strings.Split(domainInput, ",")

	// Clear allRecords before performing a new search to avoid accumulating old results
	allRecords = []DNSRecord{}

	for _, domain := range domains {
		records, err := LookupDNS(strings.TrimSpace(domain))
		if err != nil {
			http.Error(w, "Failed to lookup DNS", http.StatusInternalServerError)
			return
		}
		allRecords = append(allRecords, records...)
	}

	tmpl.ExecuteTemplate(w, "dns_results.html", allRecords)
}

func mod(x, y int) int {
	return x % y
}

// ExportHandler handles exporting allRecords to a CSV file
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	// Define the file path for the CSV file
	filePath := "dns_records.csv"

	// Export the DNS records to the CSV file
	err := ExportToCSV(allRecords, filePath)
	if err != nil {
		http.Error(w, "Failed to export CSV", http.StatusInternalServerError)
		return
	}

	// Set headers to prompt download
	w.Header().Set("Content-Disposition", "attachment; filename=dns_records.csv")
	w.Header().Set("Content-Type", "text/csv")

	// Serve the CSV file
	http.ServeFile(w, r, filePath)

	// Remove the file after serving it
	err = os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to delete file: %v", err)
	}
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	// Register custom functions
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"mod": mod,
	}).ParseGlob("templates/*.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle the DNS lookup
	http.HandleFunc("/dns-lookup", func(w http.ResponseWriter, r *http.Request) {
		DNSLookupHandler(w, r, tmpl) // Pass the parsed template
	})
	http.HandleFunc("/export", ExportHandler)

	// Serve the main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	// Start the server
	fmt.Printf("Server running at http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
