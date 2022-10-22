package node

import (
	"context"
	"log"
	"net"
	"strconv"
	"strings"

	proto "src/service"
	util "src/utilities"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const bitStrLength = 256
var prevCommitment string
var PrevAliceDiceRoll, PrevBobDiceRoll int

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

// Generate a random dice roll for Alice, and globally assign this value.
func AliceDiceRoll() int {
	PrevAliceDiceRoll = util.GenerateRandDiceRoll()
	return PrevAliceDiceRoll
}

// Generate a random dice roll for Bob, and globally assign this value.
func BobDiceRoll() int {
	PrevBobDiceRoll = util.GenerateRandDiceRoll()
	return PrevBobDiceRoll
}

// Generating the specific request for a specific node that invokes the identical method
// signature on the serverside. Here, the the response is received.
func (n *Node) InitiateRequests(c proto.ServiceClient) {
	
	// Randomly rolled dice
	log.Printf("%s rolled dice is %v", n.Name, PrevAliceDiceRoll)

	// Random k-bit string
	randStr := util.GenerateRandBitStr(bitStrLength)
	
	// Create commitment
	commCat := util.ConcatStrings(strconv.Itoa(PrevAliceDiceRoll), randStr)

	// Hashing (encoding) the concatinated k-bit string and the rolled dice.
	hashComm := util.HashStr(commCat)

	// First RPC request which includes the name and commitment.
	comm := proto.Commitment{
		Name:       n.Name,
		Commitment: hashComm,
	}

	// Call the RollDice function with the respective commitment.
	log.Printf("Sending commitment")
	diceResponse, err := c.RollDice(context.Background(), &comm)
	if err != nil {
		log.Fatalf("Error on server :: %s", err)
	}

	log.Printf("Received %s's dice roll of %v", diceResponse.Name, diceResponse.DiceRoll)

	// Second RPC request with the raw values.
	commKeys := proto.CommitmentKeys{
		Name:          n.Name, 
		DiceRoll:      int32(PrevAliceDiceRoll),
		RandomKBitStr: randStr,
	}

	// Call the CheckKeys function with the commitment key struct.
	log.Printf("Sending commitment keys")
	keyResponse, err := c.CheckKeys(context.Background(), &commKeys)
	if err != nil {
		log.Fatalf("Error on server :: %v", err)
	}

	if keyResponse.Approved {
		log.Printf("%s successfully received and validated commitment keys", keyResponse.Name)
	} else {
		log.Printf("%s received, but did not validate commitment keys", keyResponse.Name)
	}
}

// Remove Procedure Call (RPC) that returns a DiceResponse.
func (n *Node) RollDice(ctx context.Context, req *proto.Commitment) (*proto.DiceResponse, error) {
	log.Printf("Received %s's commitment", req.Name)
	prevCommitment = req.Commitment

	diceR := proto.DiceResponse{
		Name:     n.Name,
		DiceRoll: int32(PrevBobDiceRoll),
	}

	log.Printf("Sending %s's response of rolled dice %v", n.Name, PrevBobDiceRoll)
	return &diceR, nil
}

// Remove Procedure Call that returns a CorrectKeysResponse.
func (n *Node) CheckKeys(ctx context.Context, req *proto.CommitmentKeys) (*proto.CorrectKeysResponse, error) {
	log.Printf("Received %s's commitment keys", req.Name)

	var approved bool
	changeTypeDiceRoll := strconv.Itoa(int(req.DiceRoll))
	concat := util.ConcatStrings(changeTypeDiceRoll, req.RandomKBitStr)
	hash := util.HashStr(concat)
	if prevCommitment == hash {
		approved = true
	}

	keysR := proto.CorrectKeysResponse{
		Name:     n.Name,
		Approved: approved,
	}

	if approved {
		log.Printf("Sending response (acknowledgement)")
	} else {
		log.Fatalf("Sending response (denial)")
	}

	return &keysR, nil
}