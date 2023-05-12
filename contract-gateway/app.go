package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const (
	channelName   = "ivschannel"		// 連接的channel
	chaincodeName = "ivs_basic"			// 連接的chaincode
)

func main() {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	fmt.Println("getAllAssets:")
	getAllAssets(contract)

}
func getAllAssets(contract *client.Contract) {
	fmt.Println("Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("SelectAll")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func getAllUsers() ([]*User, error) {
  // Get the list of all assets.
  assets, err := GetAllAssets()
  if err != nil {
    return nil, err
  }

  // Create a slice of users.
  users := make([]*User, 0)

  // Iterate over the list of assets and create a user for each asset.
  for _, asset := range assets {
    user := &User{
      ID:   asset.ID,
      Table: asset.Table,
      Name: asset.Name,
    }
    users = append(users, user)
  }

  // Return the list of users.
  return users, nil
}
users, err := getAllUsers()
if err != nil {
  log.Fatal(err)
}

// Print the list of users.
for _, user := range users {
  fmt.Println(user)
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
