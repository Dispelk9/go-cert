package certcheck

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
)

// FetchCertWithStartTLS fetches the certificate from an SMTP server using STARTTLS.
func FetchCertWithStartTLS(host string) (string, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, "25"))
	if err != nil {
		return "", fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Perform SMTP handshake
	if err := smtpHandshake(conn, host); err != nil {
		return "", fmt.Errorf("SMTP handshake failed: %v", err)
	}

	// Upgrade to TLS
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	})
	if err := tlsConn.Handshake(); err != nil {
		return "", fmt.Errorf("TLS handshake failed: %v", err)
	}
	defer tlsConn.Close()

	// Fetch and print certificate
	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) > 0 {
		return formatCertificateDetails(state.PeerCertificates[0], host), nil
	}
	return "No peer certificates found.\n", nil
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

// Format certificate details as a string
func formatCertificateDetails(cert *x509.Certificate, host string) string {
	return fmt.Sprintf("Certificate details for %s:\n- Subject: %s\n- Issuer: %s\n- Valid From: %s\n- Valid To: %s\n- DNS Names: %v\n",
		host, cert.Subject, cert.Issuer, cert.NotBefore, cert.NotAfter, cert.DNSNames)
}

