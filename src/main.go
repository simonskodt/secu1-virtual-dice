package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"src/node"
	proto "src/service"
	util "src/utilities"

	"google.golang.org/grpc/credentials"
)

// Parse flags and commandline arguments. Create the node struct, and then setup the
// server and connect to the peers in the network. Lastly, do the RPC call in order
// to allow mutually distrustful parties (Alice and Bob) to commonly roll a dice.
func main() {

	// Flag indicating who is going to start the communication.
	// Only one of the nodes in the network should have this flag set to true.
	var initiateDiceRoll = flag.Bool("init", false, "The node that starts the die throw")
	flag.Parse()
	var offset int
	if (*initiateDiceRoll) {
		offset = 2
	}
	
	// Parse commandline arguments for a nodes port, the other notes ports,
	// and the name of the node.
	args := os.Args
	port := args[1+offset]
	portOfOtherPeers := args[2+offset]
	name := strings.Trim(args[3+offset], "`")

	// Log which node is initiated, and let the current goroutine pause for a second.
	println(util.Green + "### SETUP " + strings.ToUpper(name) + util.Reset)
	log.Printf("NODE INITIATED at port %v", strings.Trim(port, ":"))
	time.Sleep(1 * time.Second)

	// Create an instance of the Node struct which is given a name and a map to
	// the other connections of clients (including a nodes internal client).
	// A node has both a client and a server internally in order to be able
	// to create a peer-to-peer network.
	n := node.Node{
		Name:       name,
		ClientConn: make(map[string]proto.ServiceClient),
	}

	// Setting up TLS protocol. This ensures that the channel is secure.
	tlsCredentials := credentials.NewTLS(setupTLSProtocol())
	
	go n.ServerSetup(port, tlsCredentials)
	n.ConnectToPeer(portOfOtherPeers, tlsCredentials)

	// Generate Alice's and Bob's random dice rolls.
	a := node.AliceDiceRoll()
	b := node.BobDiceRoll()

	// Entry point for the different RPC calls.
	println(util.Green + "### ROLLING DICE" + util.Reset)
	rand.Seed(time.Now().UnixMicro()) // ensure that go randomization becomes non-deterministic
	for i := range n.ClientConn {
		conn := n.ClientConn[i]
		if (*initiateDiceRoll) {
			// Sending initial request from nodes with the flag -init (initiateDiceRoll) set to true.
			n.InitiateRequests(conn)
		}
	}

	time.Sleep(4 * time.Second) // have the final result appear at the same time for Alice and Bob
	println(util.Green + "\n### FINAL COMMON DICE RESULT" + util.Reset)
	log.Printf("Result: %v", util.ExclusiveOrOnTwoDiceResultsMod6Plus1(a, b))
	
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