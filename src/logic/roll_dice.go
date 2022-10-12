package logic

import (
	"context"
	"log"

	proto "src/service"
)

// Generating the specific request for a specific node that invokes the identical method
// signature on the serverside. Here, the the response is received.
func InitiateRequest(c proto.ServiceClient) {
	diceRequest := proto.DiceRequest{
		PublicKey: 55,
		Message:   72,
	}

	response, err := c.RollDice(context.Background(), &diceRequest)
	if (err != nil) {
		log.Fatalf("Server crashed :: %s", err)
	}

	log.Printf("The dice rolled a %v from %v", response.DiceOutcome, response.PublicKey)
}