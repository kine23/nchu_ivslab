package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/chaincode"
)

func main() {
	ivsChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})

	if err != nil {
		log.Panicf("Error creating ivs-transfer-basic chaincode: %v", err)
	}

	if err := ivsChaincode.Start(); err != nil {
		log.Panicf("Error starting ivs-transfer-basic chaincode: %v", err)
	}
}
