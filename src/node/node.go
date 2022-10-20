package node

import (
	"context"
	"log"
	"net"
	"strings"

	proto "src/service"
	util "src/utilities"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const bitStrLength = 256

// General accessible information on a node struct: name and the other client connections.
type Node struct {
	proto.UnimplementedServiceServer
	Name 	     string
	ClientConn   map[string]proto.ServiceClient
}

// Setting up node servers at a given port.
func (n *Node) ServerSetup(port string, cred credentials.TransportCredentials) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Server failed to listen af port %s :: %v", port, err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(cred))
	proto.RegisterServiceServer(grpcServer, n)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Servers failed to serve :: %v", err)
	}
}

// Connecting/dialing the given server ports, thus, creating a node client connection.
func (n *Node) ConnectToPeer(portOfOtherPeers string, cred credentials.TransportCredentials) {
	ports := strings.Split(portOfOtherPeers, ",")

	for i := 0; i < len(ports); i++ {
		log.Printf("Connecting to peer %s\n\n", ports[i])

		conn, err := grpc.Dial("localhost:" + ports[i], grpc.WithTransportCredentials(cred))
		if err != nil {
			log.Fatalf("Error when dialing :: %s", err)
		}

		n.ClientConn[ports[i]] = proto.NewServiceClient(conn)
	}
}

// Generating the specific request for a specific node that invokes the identical method
// signature on the serverside. Here, the the response is received.
func (n *Node) InitiateRequest(c proto.ServiceClient) {
	
	// Roll dice, resulting in a random number.
	diceRoll := util.GenerateRandDiceRoll()
	log.Printf("%s rolled dice is %v", n.Name, diceRoll)

	// Random k-bit string
	randStr := util.GenerateRandBitStr(bitStrLength)
	
	// Create commitment
	commCat := util.ConcatStrings(util.FormatInt(diceRoll), randStr)

	// Hashing (encoding) of k-bit string and the rolled dice.
	// TODO
	
	comm := proto.Commitment{
		Name:       n.Name,
		Commitment: commCat,
	}

	response, err := c.RollDice(context.Background(), &comm)
	if (err != nil) {
		log.Fatalf("Server crashed :: %s", err)
	}

	log.Printf("Received: %v", response.Received)
}

// The Remove Procedure Call (RPC) implementation that receives the request from a
// node in order to send a response back the nodes clientside.
func (n *Node) RollDice(ctx context.Context, request *proto.Commitment) (*proto.DiceResponse, error) {
	diceR := proto.DiceResponse{
		Received: true,
	}

	log.Printf("The commitment from %s is %v", request.Name, request.Commitment)

	return &diceR, nil
}