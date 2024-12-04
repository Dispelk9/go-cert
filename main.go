package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
        "github.com/Dispelk9/go-cert/certcheck" // Import the certcheck package
)

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve the home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Handle fetch-cert endpoint using the fetchCertHandler function
        http.HandleFunc("/fetch-cert", fetchCertHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// Home page handler to serve the HTML file
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func fetchCertHandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error looking up MX records: %v", err), http.StatusInternalServerError)
		return
	}

	if len(mxRecords) == 0 {
		http.Error(w, "No MX records found", http.StatusNotFound)
		return
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("MX records for %s:\n", domain))
	for _, mx := range mxRecords {
		response.WriteString(fmt.Sprintf("- %s (Priority: %d)\n", mx.Host, mx.Pref))
	}

	for _, mx := range mxRecords {
		mxHost := strings.TrimSuffix(mx.Host, ".")
		response.WriteString(fmt.Sprintf("\nFetching certificate from %s...\n", mxHost))

		certDetails, err := certcheck.FetchCertWithStartTLS(mxHost)
		if err != nil {
			response.WriteString(fmt.Sprintf("Failed to fetch certificate from %s: %v\n", mxHost, err))
			continue
		}

		response.WriteString(certDetails + "\n")
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response.String()))
}



