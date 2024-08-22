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


//handler

func DNSLookupHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	domainInput := r.FormValue("domains")
	domains := strings.Split(domainInput, ",")

	var allRecords []DNSRecord
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

func main() {

	//load envfile
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

	// Serve the main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	// Start the server
	fmt.Printf(`server running http://localhost%s`, port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

