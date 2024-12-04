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

	fmt.Println("Server is running on http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}


// Home page handler to serve the HTML file
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func fetchCertHandler(w http.ResponseWriter, r *http.Request) {
    domains := r.URL.Query().Get("domain") // Fetch input for domains
    if domains == "" {
        http.Error(w, "Domain is required", http.StatusBadRequest)
        return
    }

    // Split the input using multiple delimiters: commas, spaces, and newlines
    separators := strings.NewReplacer(",", " ", "\n", " ")
    normalizedDomains := separators.Replace(domains)
    domainList := strings.Fields(normalizedDomains) // Split by whitespace

    var response strings.Builder

    for _, domain := range domainList {
        domain = strings.TrimSpace(domain) // Clean up whitespace
        if domain == "" {
            continue
        }

        response.WriteString(fmt.Sprintf("Processing domain: %s\n", domain))

        mxRecords, err := net.LookupMX(domain)
        if err != nil {
            response.WriteString(fmt.Sprintf("Error looking up MX records for domain %s: %v\n", domain, err))
            continue
        }

        if len(mxRecords) == 0 {
            response.WriteString(fmt.Sprintf("No MX records found for domain %s\n", domain))
            continue
        }

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

        response.WriteString("\n---\n") // Separator for each domain
    }

    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(response.String()))
}

