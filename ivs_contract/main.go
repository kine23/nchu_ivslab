package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/contract"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&contract.UserContract{}, &contract.IVSContract{})
	if err != nil {
		panic(err)
	}

	if err := chaincode.Start(); err != nil {
		panic(err)
	}
}
