package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/rajan-31/hyperledger-fabric_chaincode-devmode-contractapi/my_chaincode"
)

func main() {

	chaincode, err := contractapi.NewChaincode(new(my_chaincode.SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
