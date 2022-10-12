package node

import (
	"context"
	"log"
	"net"
	"strings"

	proto "src/service"
	"google.golang.org/grpc"
)

type Node struct {
	proto.UnimplementedServiceServer
	Name 	     string
	ClientConn   map[string]proto.ServiceClient
}

// Setting up node servers at a given port.
func (n *Node) ServerSetup(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Server failed to listen af port %s :: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterServiceServer(grpcServer, n)

	log.Printf("Server listens at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Servers failed to serve :: %v", err)
	}
}

// Connecting/dialing the given server ports, thus, creating a node client connection.
func (n *Node) ConnectToPeer(portOfOtherPeers string) {
	ports := strings.Split(portOfOtherPeers, ",")
	
	for i := 0; i < len(ports); i++ {
		log.Printf("Connecting to peer %s", ports[i])
		conn, err := grpc.Dial("localhost:" + ports[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Error when dialing :: %s", err)
		}
		n.ClientConn[ports[i]] = proto.NewServiceClient(conn)
	}
}

// The Remove Procedure Call (RPC) implementation that receives the request from a
// node in order to send a response back the nodes clientside.
func (n *Node) RollDice(ctx context.Context, request *proto.DiceRequest) (*proto.DiceResponse, error) {
	diceR := proto.DiceResponse{
		PublicKey: request.PublicKey + 1000,
		DiceOutcome: request.Message - 1000,
	}

	log.Printf("The given request is %v", request.Message)

	return &diceR, nil
}