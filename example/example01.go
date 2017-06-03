/* 
    chaincode for sharing history storage

    author: Runshan Hu
*/


package main

import (
        "errors"
        "fmt"
        "github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {

}

func main() {
        err := shim.Start(new(SimpleChaincode))
        if err != nil {
                fmt.Printf("Error starting sharing historty storage chaincode: %s", err)
        }
}

//Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
        if len(args) != 2 {
                return nil, errors.New("Incorrect number of arguments. Expecting 2")
        }

        err := stub.PutState(args[0], []byte(args[1]))
        if err != nil {
                return nil, err
        }
        
        return nil, nil
}

//Invoke entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
        fmt.Println("invoke is running" + function)

        //Handle different functions
        if function == "init" {
                return t.Init(stub, "init", args)
        } else if function == "write" {
                return t.write(stub, args)
        }

        fmt.Println("invoke did not find func:" + function)

        return nil, errors.New("Received unknown function invocation")

}

//Query is entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
        var key, jsonResp string
        var err error

        if len(args) != 1 {
                return nil, errors.New("Incorrect number of arguments. Expecting name of the dataset to query")
        }

        key = args[0]
        valAsbytes, err := stub.GetState(key)

        if err != nil {
                jsonResp = "{\"Error\": \"Failed to get the state for " + key + "\"}"
                return nil, errors.New(jsonResp)
        }

        return valAsbytes, nil
}

//write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
        var datasetId, value string
        var err error
        fmt.Println("running write()")

        if len(args) != 2 {
                return nil, errors.New("Incorrect number of arguments. Expecting 2. DatasetID and value to set")
        }

        datasetId = args[0]
        value = args[1]
        
        //write the variable into the chaincode state
        err = stub.PutState(datasetId, []byte(value))
        if err != nil {
                return nil, err
        }

        return nil, nil
}

/*
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
        var key, jsonResp string
        var err error

        if len(args) !=1 {
                return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
        }

        key = args[0]
        valAsbytes, err := stub.GetState(key)

        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
                return nil, errors.New(jsonResp)
        }
        
        return valAsbytes, nil
}
*/


