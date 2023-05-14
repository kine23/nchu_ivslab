package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/kine23/nchu_ivslab/ivs_contract/model"
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

	// Create a new project
	var project model.Project
	fmt.Println("Enter Manufacturer:")
	fmt.Scanln(&project.Manufacturer)
	fmt.Println("Enter ManufactureLocation:")
	fmt.Scanln(&project.ManufactureLocation)
	fmt.Println("Enter PartName:")
	fmt.Scanln(&project.PartName)
	fmt.Println("Enter BatchNumber:")
	fmt.Scanln(&project.BatchNumber)	
	fmt.Println("Enter SerialNumber:")
	fmt.Scanln(&project.SerialNumber)
	fmt.Println("Enter ManufactureDate:")
	fmt.Scanln(&project.ManufactureDate)
	fmt.Println("Enter Organization:")
	fmt.Scanln(&project.Organization)
	// 寫入新零件
	fmt.Println("createParts:")
	createParts(contract, project)
	// 轉移零件
	var projectID string
	var newOwnerOrganization string
	fmt.Println("Enter the ID of the project to transfer:")
	fmt.Scanln(&projectID)
	fmt.Println("Enter the new owner organization:")
	fmt.Scanln(&newOwnerOrganization)
	fmt.Println("transferProject:")
	transferProject(contract, projectID, newOwnerOrganization)
	// 刪除零件
	fmt.Println("Enter the ID of the project to delete:")
	fmt.Scanln(&projectID)
	fmt.Println("deleteProject:")
	deleteProject(contract, projectID)
	// 查詢指定項目
	fmt.Println("Enter the ID of the project to select:")
	fmt.Scanln(&projectID)
	fmt.Println("selectProject:")
	selectProject(contract, projectID)
	// 查詢索引項目
	var key string
	var value string
	fmt.Println("Enter the key to select by:")
	fmt.Scanln(&key)
	fmt.Println("Enter the value to select by:")
	fmt.Scanln(&value)
	fmt.Println("selectBySome:")
	selectBySome(contract, key, value)
	// 查詢全部項目
	fmt.Println("getAllParts:")
	getAllParts(contract)
	// 多頁查詢項目
	var pageSize int32
	var bookmark string
	fmt.Println("Enter the page size for pagination:")
	fmt.Scanln(&pageSize)
	fmt.Println("Enter the bookmark for pagination:")
	fmt.Scanln(&bookmark)
	fmt.Println("selectAllWithPagination:")
	selectAllWithPagination(contract, pageSize, bookmark)
	// 註冊新使用者
	var username string
	var name string
	fmt.Println("Enter username:")
	fmt.Scanln(&username)
	fmt.Println("Enter name:")
	fmt.Scanln(&name)
	fmt.Println("createUser:")
	createUser(contract, username, name)
	// 讀取使用者資訊
	var username string
	fmt.Println("Enter username to read:")
	fmt.Scanln(&username)
	fmt.Println("readUser:")
	readUser(contract, username)
	// 查詢全部User
	fmt.Println("getAllUsers:")
	getAllUsers(contract)
}

// 寫入新零件
func createParts(contract *client.Contract, project model.Project) {
	fmt.Println("Submit Transaction: Create, function inserts a new Parts into the ledger")

	projectJSON, err := json.Marshal(project)
	if err != nil {
		panic(fmt.Errorf("failed to marshal project: %w", err))
	}

	submitResult, err := contract.SubmitTransaction("Insert", string(projectJSON))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}
	result := formatJSON(submitResult)

	fmt.Printf("*** Result: %s\n", result)
}

// 轉移零件
func transferProject(contract *client.Contract, projectID string, newOwnerOrganization string) {
	fmt.Println("Submit Transaction: TransferProject, function transfers the ownership of a project to a new organization")

	submitResult, err := contract.SubmitTransaction("TransferProject", projectID, newOwnerOrganization)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}
	result := formatJSON(submitResult)

	fmt.Printf("*** Result: %s\n", result)
}

// 刪除零件
func deleteProject(contract *client.Contract, projectID string) {
	fmt.Println("Submit Transaction: Delete, function deletes a project from the ledger")

	submitResult, err := contract.SubmitTransaction("Delete", projectID)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}
	result := formatJSON(submitResult)

	fmt.Printf("*** Result: %s\n", result)
}

// 查詢指定項目
func selectProject(contract *client.Contract, projectID string) {
	fmt.Println("Evaluate Transaction: SelectByIndex, function retrieves a project from the ledger")

	evaluateResult, err := contract.EvaluateTransaction("SelectByIndex", projectID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result: %s\n", result)
}

// 查詢索引項目
func selectBySome(contract *client.Contract, key string, value string) {
	fmt.Println("Evaluate Transaction: SelectBySome, function retrieves projects from the ledger by some key and value")

	evaluateResult, err := contract.EvaluateTransaction("SelectBySome", key, value)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result: %s\n", result)
}

// 查詢全部項目
func getAllParts(contract *client.Contract) {
	fmt.Println("Submit Transaction: GetAllParts, function returns all the current assets on the ledger")

	evaluateResult, err := contract.SubmitTransaction("SelectAll")
	if err != nil {
		panic(fmt.Errorf("failed to Submit transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// 多頁查詢項目
func selectAllWithPagination(contract *client.Contract, pageSize int32, bookmark string) {
	fmt.Println("Evaluate Transaction: SelectAllWithPagination, function retrieves all projects from the ledger with pagination")

	evaluateResult, err := contract.EvaluateTransaction("SelectAllWithPagination", fmt.Sprintf("%d", pageSize), bookmark)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result: %s\n", result)
}
// 註冊新使用者
func createUser(contract *client.Contract, username string, name string) {
	fmt.Println("Submit Transaction: CreateUser, function creates a new user")

	submitResult, err := contract.SubmitTransaction("CreateUser", username, name)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}
	result := formatJSON(submitResult)

	fmt.Printf("*** Result: %s\n", result)
}

// 讀取使用者資訊
func readUser(contract *client.Contract, username string) {
	fmt.Println("Evaluate Transaction: ReadUser, function returns the details of a user")

	evaluateResult, err := contract.EvaluateTransaction("ReadUser", username)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result: %s\n", result)
}
	
// 查詢全部User
func getAllUsers(contract *client.Contract) {
fmt.Println("Submit Transaction: GetAllUsers, function returns all the current users on the ledger")
	evaluateResult, err := contract.SubmitTransaction("GetAllUsers")
if err != nil {
	panic(fmt.Errorf("failed to Submit transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}


func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
