package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const (
	channelName   = "p-channel1"		// 連接的channel
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
	
// Project項目列表
type Project struct {
	Table               string `json:"table" form:"table"`  //數據庫標記
	Manufacturer        string `json:"Manufacturer"`        //製造商
	ManufactureLocation string `json:"ManufactureLocation"` //製造地點
	PartName            string `json:"PartName"`            //零件名稱
	PartNumber          string `json:"PartNumber"`          //零件批號
	SerialNumber        string `json:"SerialNumber"`        //產品序號
	ManufactureDate     string `json:"ManufactureDate"`     //製造日期
	Item                string `json:"Item"`                //項目
	ID                  string `json:"ID"`                  //項目唯一ID
	Category            string `json:"Category"`            //所屬類別
	Describes           string `json:"Describes"`           //描述
	Developer           string `json:"Developer"`           //開發者
	Organization        string `json:"Organization"`        //組織
}
	
	project := Project{
		Manufacturer:        "Serurity",
		ManufactureLocation: "Taiwan",
		PartName:            "Serurity-Chip-V1",
		BatchNumber:         "PNSPAQ11230098",
		SerialNumber:        "SNSSAQ11230098",
		ManufactureDate:     "2023-05-14",
		Organization:        "Security-Org",
	}
	insertProject(contract, project)
//	transferProject(contract, project.ID, "Brang-Org")
//	deleteProject(contract, project.ID)
//	getProjectBySerialNumber(contract, project.SerialNumber)
//	fmt.Println("getAllAssets:")
//	getAllAssets(contract)
//	fmt.Println("getAllUsers:")
//	getAllUsers(contract)
}
func insertProject(contract *client.Contract, project Project) { 
	fmt.Println("Submit Transaction: Insert, function inserts a new project into the ledger")

	projectJSON, err := json.Marshal(project)
	if err != nil {
		panic(fmt.Errorf("failed to marshal project: %w", err))
	}

	submitResult, err := contract.SubmitTransaction("Insert", string(projectJSON))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	result := formatJSON(submitResult)
	fmt.Printf("*** Result:%s\n", result)
}

//func transferProject(contract *client.Contract, projectID string, newOwnerOrganization string) {
//	fmt.Println("Submit Transaction: TransferProject, function transfers the ownership of a project to a new organization")
//
//	submitResult, err := contract.SubmitTransaction("TransferProject", projectID, newOwnerOrganization)
//	if err != nil {
//		panic(fmt.Errorf("failed to submit transaction: %w", err))
//	}
//
//	result := formatJSON(submitResult)
//
//	fmt.Printf("*** Result:%s\n", result)
//}

//func deleteProject(contract *client.Contract, projectID string) {
//	fmt.Println("Submit Transaction: Delete, function deletes a project from the ledger")
//
//	submitResult, err := contract.SubmitTransaction("Delete", projectID)
//	if err != nil {
//		panic(fmt.Errorf("failed to submit transaction: %w", err))
//	}
//
//	result := formatJSON(submitResult)
//
//	fmt.Printf("*** Result:%s\n", result)
//}

//func getProjectBySerialNumber(contract *client.Contract, serialNumber string) {
//	fmt.Println("Evaluate Transaction: SelectBySome, function returns a project with the specified serial number from the ledger")
//
//	evaluateResult, err := contract.EvaluateTransaction("SelectBySome", "SerialNumber", SNSSAQ11230098)
//	if err != nil {
//		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
//	}
//
//	result := formatJSON(evaluateResult)
//
//	fmt.Printf("*** Result:%s\n", result)
//}

//func getAllAssets(contract *client.Contract) {
//	fmt.Println("Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
//
//	evaluateResult, err := contract.EvaluateTransaction("SelectAll")
//	if err != nil {
//		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
//	}
//	result := formatJSON(evaluateResult)
//
//	fmt.Printf("*** Result:%s\n", result)
//}

//func getAllUsers(contract *client.Contract) {
//fmt.Println("Evaluate Transaction: GetAllUsers, function returns all the current users on the ledger")
//	evaluateResult, err := contract.EvaluateTransaction("GetAllUsers")

//	if err != nil {
//	panic(fmt.Errorf("failed to evaluate transaction: %w", err))
//	}
//	result := formatJSON(evaluateResult)
//
//	fmt.Printf("*** Result:%s\n", result)
//}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
