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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/fetch-cert", fetchCertHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Home page handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>SMTP Certificate Fetcher</title>
		</head>
		<body>
			<h1>Fetch SMTP Certificate</h1>
			<form action="/fetch-cert" method="get">
				<label for="domain">Domain:</label>
				<input type="text" id="domain" name="domain" required>
				<br>
				<button type="submit">Fetch Certificate</button>
			</form>
		</body>
		</html>
	`)
}

// Fetch certificate handler
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

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "MX records for %s:\n", domain)
	for _, mx := range mxRecords {
		fmt.Fprintf(w, "- %s (Priority: %d)\n", mx.Host, mx.Pref)
	}

	for _, mx := range mxRecords {
		mxHost := strings.TrimSuffix(mx.Host, ".")
		fmt.Fprintf(w, "\nFetching certificate from %s...\n", mxHost)

		certDetails, err := certcheck.FetchCertWithStartTLS(mxHost)
		if err != nil {
			fmt.Fprintf(w, "Failed to fetch certificate from %s: %v\n", mxHost, err)
			continue
		}

		fmt.Fprintln(w, certDetails)
	}
}

