package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Function to print certificate details
func printCertificateDetails(cert *x509.Certificate, host string) {
	fmt.Printf("Certificate details for %s:\n", host)
	fmt.Printf("- Subject: %s\n", cert.Subject)
	fmt.Printf("- Issuer: %s\n", cert.Issuer)
	fmt.Printf("- Valid From: %s\n", cert.NotBefore)
	fmt.Printf("- Valid To: %s\n", cert.NotAfter)
	fmt.Printf("- DNS Names: %v\n", cert.DNSNames)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <domain>")
		return
	}

	domain := os.Args[1]

	// Step 1: Get MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Fatalf("Error looking up MX records for domain %s: %v", domain, err)
	}

	if len(mxRecords) == 0 {
		log.Fatalf("No MX records found for domain %s", domain)
	}

	fmt.Printf("MX records for %s:\n", domain)
	for _, mx := range mxRecords {
		fmt.Printf("- %s (Priority: %d)\n", mx.Host, mx.Pref)
	}

	// Step 2: Fetch certificates from each MX server
	for _, mx := range mxRecords {
		mxHost := strings.TrimSuffix(mx.Host, ".") // Trim trailing dot
		fmt.Printf("\nFetching certificate from %s...\n", mxHost)

		err := fetchCertWithStartTLS(mxHost)
		if err != nil {
			log.Printf("Failed to fetch certificate from %s: %v\n", mxHost, err)
		}
	}
}

// Fetch certificate using STARTTLS
func fetchCertWithStartTLS(host string) error {
	// Connect to the MX server on port 25
	conn, err := net.Dial("tcp", net.JoinHostPort(host, "25"))
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Perform SMTP handshake
	if err := smtpHandshake(conn, host); err != nil {
		return fmt.Errorf("SMTP handshake failed: %v", err)
	}

	// Upgrade to TLS
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	})
	if err := tlsConn.Handshake(); err != nil {
		return fmt.Errorf("TLS handshake failed: %v", err)
	}
	defer tlsConn.Close()

	// Fetch and print certificate
	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) > 0 {
		printCertificateDetails(state.PeerCertificates[0], host)
	} else {
		fmt.Println("No peer certificates found.")
	}
	return nil
}

// Perform the SMTP handshake and issue STARTTLS
func smtpHandshake(conn net.Conn, host string) error {
	buf := make([]byte, 1024)

	// Read the server's initial response
	if _, err := conn.Read(buf); err != nil {
		return fmt.Errorf("failed to read initial response: %v", err)
	}

	// Send EHLO command
	ehloCmd := fmt.Sprintf("EHLO %s\r\n", host)
	if _, err := conn.Write([]byte(ehloCmd)); err != nil {
		return fmt.Errorf("failed to send EHLO command: %v", err)
	}

	// Read EHLO response
	if _, err := conn.Read(buf); err != nil {
		return fmt.Errorf("failed to read EHLO response: %v", err)
	}

	// Send STARTTLS command
	starttlsCmd := "STARTTLS\r\n"
	if _, err := conn.Write([]byte(starttlsCmd)); err != nil {
		return fmt.Errorf("failed to send STARTTLS command: %v", err)
	}

	// Read STARTTLS response
	if _, err := conn.Read(buf); err != nil {
		return fmt.Errorf("failed to read STARTTLS response: %v", err)
	}

	return nil
}

