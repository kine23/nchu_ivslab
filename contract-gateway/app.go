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
	chaincodeName = "ivs_basic"		// 連接的chaincode
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
	fmt.Println("getAllUsers:")
	getAllUsers(contract)
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

func getAllUsers(contract *client.Contract) {
	fmt.Println("Evaluate Transaction: GetAllUsers, function returns all the current users on the ledger")
	evaluateResult, err := contract.EvaluateTransaction("GetAllUsers")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	// 解析評估結果為使用者清單
	var users []*User
	if err := json.Unmarshal(evaluateResult, &users); err != nil {
		panic(fmt.Errorf("failed to unmarshal users: %w", err))
	}

	// 輸出所有使用者清單
	result, err := json.Marshal(users)
	if err != nil {
		panic(fmt.Errorf("failed to marshal users: %w", err))
	}
	fmt.Printf("*** Result:%s\n", string(result))
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
