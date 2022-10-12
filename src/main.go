package main

import (
	"log"
	"os"

	"src/logic"
	"src/node"
	proto "src/service"
)

func main() {
	log.Println("NODE INITIATED")
	// fmt.Printf("Format of input :: <:port1> <port2,port3> <name>")
	// fmt.Printf("        Example :: :9000 9001 Alice")
	// fmt.Printf("        Example :: :9001 9002 Bob")

	args := os.Args
	port := args[1]
	portOfOtherPeers := args[2]
	name := args[3]

	n := node.Node{
		Name:       name,
		ClientConn: make(map[string]proto.ServiceClient),
	}
	
	go n.ServerSetup(port)
	n.ConnectToPeer(portOfOtherPeers)
	for i := range n.ClientConn {
		conn := n.ClientConn[i]
		logic.InitiateRequest(conn)
	}
	
	for {} // Prevent termination
}