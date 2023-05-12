package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/contract"
)

func main() {
	IVSContract := new(IVSContract)
	IVSContract.Contract = new(contractapi.Contract)

	userContract := new(UserContract)
	userContract.Contract = new(contractapi.Contract)

	chaincode, err := contractapi.NewChaincode(IVSContract, userContract)
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}

}

//func main() {
//	chaincode, err := contractapi.NewChaincode(&contract.UserContract{}, &contract.IVSContract{})
//	if err != nil {
//		panic(err)
//	}
//
//	if err := chaincode.Start(); err != nil {
//		panic(err)
//	}
//}
