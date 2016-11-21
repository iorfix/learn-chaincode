/*
Copyright IBM Corp 2016 All Rights Reserved.

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
	"strconv"
	"time"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Waste struct {
	Id					string	`json:"id"` 
	Producer        	string	`json:"producer"`
	QuantityProduced    int		`json:"quantityProduced"`
	TimestampProduced	int64	`json:"timestampProduced"`	//utc timestamp of creation
	TimestampAssigned	int64	`json:"timestampAssigned"`	//utc timestamp of assignment
	Retriever			string  `json:"retriever"`
	TimestampRetrieved	int64	`json:"timestampRetrieved"`	//utc timestamp of assignment
	QualityRetrieved    int 	`json:"qualityRetrieved"`
}

type Waste_Holder struct {
	wId 	[]string `json:"wids"`
}


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var wasteIDs Waste_Holder
//	var bytes [5]byte 
	//for debug
//	wasteIDs.wId = [5]string{'1', '2', 'A', 'B', 'AA'}

	bytes, err := json.Marshal(wasteIDs)

    if err != nil { 
		return nil, errors.New("Error creating wasteIDs record") 
	}

	err = stub.PutState("wasteIDs", bytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}


// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	user, err := t.get_username(stub)

	if err != nil { return nil, errors.New("Error retrieving caller information")}

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "newWaste" {
		return t.newWaste(stub, user, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "readWaste" { //read a variable
		waste, err := t.readWasteB(stub, args)
		if err != nil { return nil, err}
		//return json.Marshal(waste)
		 return waste, err
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query...: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) newWaste(stub shim.ChaincodeStubInterface, user string, args []string) ([]byte, error) {
	var id string
	var quantity int
	var timestamp int64
	var err error
	var waste Waste
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. id, quantity")
	}

	id = args[0] //rename for funsies
	quantity, err = strconv.Atoi(args[1])
	timestamp = makeTimestamp()
	
	waste.Id = id
	waste.Producer = user
	waste.QuantityProduced = quantity
	waste.TimestampProduced = timestamp
	
	wByte, err := json.Marshal(waste)
	if err != nil {
		return nil, err
	}
	
	err = stub.PutState(id, []byte(wByte)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *SimpleChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

    username, err := stub.ReadCertAttribute("username");
	if err != nil { return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error()) }
	return string(username), nil
}


func (t *SimpleChaincode) readWasteB(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Waste id")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, err
}

// read - query function to read key/value pair
func (t *SimpleChaincode) readWaste(stub shim.ChaincodeStubInterface, args []string) (Waste, error) {
	var key, jsonResp string
	var waste Waste
	
	if len(args) != 1 {
		return waste, errors.New("Incorrect number of arguments. Expecting name of the Waste id")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return waste, errors.New(jsonResp)
	}
	err = json.Unmarshal(valAsbytes, &waste);
    if err != nil {	
		fmt.Printf("retrieve WASTE: Corrupt Waste "+string(valAsbytes)+": %s", err) 
		return waste, errors.New("RETRIEVE_WASTE: Corrupt waste record"+string(valAsbytes))
	}
	return waste, nil
}
