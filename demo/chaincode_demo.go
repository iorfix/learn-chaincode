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
	"encoding/binary"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type OpenBinObj struct {
	Id				  uint32	`json:"id"` 
	Producer          string	`json:"producer"`
	Lat				  float64	`json:"lat"`
	Lng				  float64	`json:"lng"`
	TimestampOpened	  int64	`json:"timestampOpened"`	//utc timestamp of creation
	TimestampClosed	  int64	`json:"timestampClosed"`
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
	var users []string
	usersByte, err2 := json.Marshal(users)
	if (err2 !=nil) {
		return nil, err2
	}
	err := stub.PutState("USERLIST", usersByte)
	return nil, err
}


// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	caller, er := stub.GetPayload()
	if (er !=nil) {
		fmt.Println(er)
	}
	fmt.Println("GetPayload: ", caller)
	

	//user, err := t.get_username(stub)
	//user := "PIPPO"
	
	//fmt.Println("username: " + user)
	
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "newOpening" {
		return t.newOpening(stub, args)
//	} else if function == "collect" {
//		user := "COLL"
//		return t.collectWaste(stub, user, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function) 
	fmt.Println(args)
	// Handle different functions
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting key")
	}
	key := args[0]
	fmt.Println("readkey:", key)
	var byteVal []byte
	var err error
	if (function == "readall") {
		byteVal, err = readAllFromUser(stub, key)
	} else {
		byteVal, err = readKeyState(stub, key)
	}
	
	fmt.Println("reading:", byteVal)
	return byteVal, err
}

func (t *SimpleChaincode) newOpening(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Opening:",  args)
		if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting user, lat, lng, open, close")
	}
	user := args[0]
	var chainuserarray []byte
	var err error
	var id uint32
	var openbin OpenBinObj
	chainuserarray, err = readChain (stub, user);
	if (err != nil) {
		return nil, err
	}
	if (chainuserarray == nil) {
		chainuserarray = make([]byte, 0)
		erradd := addnewuser(stub, user)
		if (erradd !=nil) {
			return nil, erradd
		}
	}
	fmt.Println("prev chainuserarray", chainuserarray )
	id = makeTimestamp()
	idByteArr := make([]byte, 4)
	idS := strconv.FormatUint(uint64(id), 10)
    binary.LittleEndian.PutUint32(idByteArr, id)
	chainuserarray = append(chainuserarray, idByteArr...)
	fmt.Println("new chainuserarray", chainuserarray)
	err = writeUserChain(stub, user, chainuserarray)
	
	openbin.Id = id
	openbin.Producer = user
	openbin.Lat, err = strconv.ParseFloat(args[1], 64)
	if (err !=nil) {
		return nil, err
	}
	openbin.Lng, err = strconv.ParseFloat(args[2], 64)
	if (err !=nil) {
		return nil, err
	}
	openbin.TimestampOpened, err = strconv.ParseInt(args[3], 10, 64)
	if (err !=nil) {
		return nil, err
	}
	openbin.TimestampClosed, err = strconv.ParseInt(args[4], 10, 64)
	if (err !=nil) {
		return nil, err
	}
	openBinByte, err2 := json.Marshal(openbin)
	if (err2 !=nil) {
		return nil, err2
	}
	fmt.Println("Marshalled:" + string(openBinByte))
	
	err = stub.PutState(idS, openBinByte)
	return nil, err
	
}

func addnewuser(stub shim.ChaincodeStubInterface, newuser string) (error) {
	valAsbytes, err := stub.GetState("USERLIST")
	if (err != nil) {
		return err
	}
	var userlist []string
	err = json.Unmarshal(valAsbytes, &userlist);
	userlist = append(userlist, newuser)
	fmt.Println("new userlist:" , userlist)
	wByte, err2 := json.Marshal(&userlist)
	if (err2!=nil) {
		return err2
	}
	return stub.PutState("USERLIST", wByte)
}

func readChain(stub shim.ChaincodeStubInterface, user string) ([]byte, error) {
	valAsbytes, err := stub.GetState(user)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	
	fmt.Println(valAsbytes)
	return valAsbytes, err
}

func  writeUserChain(stub shim.ChaincodeStubInterface, user string, vals []byte) (error) {
	err := stub.PutState(user, vals) //write the variable into the chaincode state
	return err
}


// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() uint32 {
	var now int64 
	now = time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
    return uint32(now)
}


func readKeyState(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	var jsonResp string
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println("Retrieving:" + string(valAsbytes))
	return valAsbytes, err
}

func readAllFromUser(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	
	chainuserarray, err := readChain(stub, key)
	if err != nil {
		return nil, err
	}

	chainuserint := convertByteArrayToUint32Array(&chainuserarray)
	fmt.Println(chainuserint)
	numElems := len(chainuserint)
	openBinArr := make([]OpenBinObj, numElems)
	for i := 0; i < numElems; i++ {
		idConf :=  strconv.FormatUint(uint64(chainuserint[i]),10)
		valAsbytes, err := readKeyState(stub, idConf)
		if (err !=nil) {
			fmt.Println("error reading:", idConf)
			return nil, err
		}
		err = json.Unmarshal(valAsbytes, &openBinArr[i]);
		fmt.Println("Read:", openBinArr[i])
	}
	fmt.Println("Fullread:", openBinArr)
	
	wByte, err2 := json.Marshal(&openBinArr)
	return wByte, err2
}

func convertByteArrayToUint32Array(bytearray *[]byte) []uint32 {
	numElems := len(*bytearray)/4
	uintarray := make([]uint32, numElems)
	for i := 0; i < numElems; i++ {
		elemByte := (*bytearray)[i*4:i*4+4]
		val := binary.LittleEndian.Uint32(elemByte)
		uintarray[i] = val
	}
	fmt.Println(uintarray)
	return uintarray
}




