package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"

	"src/logic"
	"src/node"
	proto "src/service"

	"google.golang.org/grpc/credentials"
)

func main() {
	args := os.Args
	port := args[1]
	portOfOtherPeers := args[2]
	name := strings.Trim(args[3], "`")

	log.Printf("NODE :: %v INITIATED", strings.ToUpper(name))

	n := node.Node{
		Name:       name,
		ClientConn: make(map[string]proto.ServiceClient),
	}

	// Setting up TLS protocol
	tlsCredentials := credentials.NewTLS(setupTLSProtocol())
	
	go n.ServerSetup(port, tlsCredentials)
	n.ConnectToPeer(portOfOtherPeers, tlsCredentials)
	for i := range n.ClientConn {
		conn := n.ClientConn[i]
		logic.InitiateRequest(conn)
	}
	
	for {} // Prevent termination
}

func setupTLSProtocol() *tls.Config {
	certX509Pool := x509.NewCertPool()
	certs := []tls.Certificate{}

	for _, nodeName := range []string{"alice", "bob"} {

		// Include certificate files
		certContent, err := os.ReadFile(fmt.Sprintf("certificates/%v.cert.pem", nodeName))
		if (err != nil) {
			log.Fatalf("Could not read certificates :: %v", err)
		}

		// Parsing of certificates
		decodeCert, _ := pem.Decode(certContent)
		clientCertificate, err := x509.ParseCertificate(decodeCert.Bytes)
		if (err != nil) {
			log.Fatalf("Could not parse the client certificate :: %v", err)
		}

		// Auth and self-signing
		clientCertificate.BasicConstraintsValid = true
		clientCertificate.IsCA = true
		clientCertificate.KeyUsage = x509.KeyUsageCertSign
		certX509Pool.AppendCertsFromPEM(certContent)

		// Loading certificates
		loadCertificates, err := tls.LoadX509KeyPair(
				fmt.Sprintf("certificates/%v.cert.pem", nodeName), 
				fmt.Sprintf("certificates/%v.key.pem", nodeName),
			)
		if (err != nil) {
			log.Fatalf("Could not load X509 keypair :: %v", err)
		}

		certs = append(certs, loadCertificates)
	}

	return &tls.Config{
		Certificates: certs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certX509Pool,
		RootCAs:      certX509Pool,
	}
}