package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"fabriclab.com/mylab_ivs/ivs_contract/contract"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&contract.UserContract{}, &contract.ProjectContract{})
	if err != nil {
		panic(err)
	}

	if err := chaincode.Start(); err != nil {
		panic(err)
	}
}
