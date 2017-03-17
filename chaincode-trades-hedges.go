/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// TradesHedgesData to store trades and hedges in one common-data model.
type TradesHedgesData struct {

	TradeID string `json:"tradeID"`
    TradeIndicator  string `json:"tradeIndicator"`
    TradeType       string `json:"tradeType"`
    CurrencyFrom     string `json:"currencyFrom"`
    CurrencyTo    string `json:"currencyTo"`
    TradeDirection       string `json:"tradeType"`
    Rate     string `json:"rate"`
    TotalTrade    string `json:"totalTrade"`    
    TradeDate     string `json:"tradeDate"`
    SettlementDate    string `json:"settlementDate"` 
}

func main() {
	err := shim.Start(new(TradesHedgesData))
	if err != nil {
		fmt.Printf("Error starting trades-hedges chaincode: %s", err)
	}
}

// Init resets all the things
func (t *TradesHedgesData) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("trades_hedges", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *TradesHedgesData) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *TradesHedgesData) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *TradesHedgesData) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	var bytes []byte
	var merr error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies

	//args[1] is json document  with the trades data for TradesHedgesData common data model
	value = args[1]

	var tradeshedgesInfo TradesHedgesData

	fmt.Println("Data being added to block for " + key + " Value: " + value)

	//Converting the string to object
	merr = json.Unmarshal([]byte(value),tradeshedgesInfo)

	//converting object to bytes for storing it into ledger
	bytes, merr = json.Marshal(tradeshedgesInfo)

	if merr != nil {
		fmt.Println("Couldn't marshan tradeshedgesInfo object")
		return nil, merr
	}

		//write the trades info into the chaincode state
		err = stub.PutState(key, bytes) 


	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *TradesHedgesData) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)


	fmt.Println("Printing retrieved response: " + string(valAsbytes))

	var tradeshedgesInfo TradesHedgesData

	
	if err != nil {
		jsonResp = "{\"Error\":\"Step 1: Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}


	//Converting the string to object
	err = json.Unmarshal(valAsbytes, &tradeshedgesInfo)

	if err != nil {
		jsonResp = "{\"Error\":\"Step 2: Failed to Unmarshal " + key + "\"}"
		return nil, errors.New(jsonResp)
	}


	fmt.Println(tradeshedgesInfo.TradeID)


	return valAsbytes, nil
}
