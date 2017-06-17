/* 
    chaincode for sharing history storage

    author: Runshan Hu
*/


package main

import (
        "errors"
        "fmt"
        "log"
        "encoding/json"
        "net/http"
        "time"
        "io/ioutil"
        "github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

type dataResult struct {
        ID             string 
        Name           string
        Salary         float64
        V              int
        CreatedDate    time.Time
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
        getResultAnonyService(); 
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

func getResultAnonyService() {
        resp, err := http.Get("http://10.7.6.25:3000/dataset/all")
        if err != nil {
                log.Println(err);
        }

        body,err := ioutil.ReadAll(resp.Body);
        if err != nil {
                log.Println(err);
        }

        defer resp.Body.Close();

        result := make([]dataResult, 0);

        if err := json.Unmarshal(body, &result); err != nil {
                log.Println(err); 
        }

        for index, element := range result {
                fmt.Println("-->Record ", index);
                fmt.Println("name = ", element.Name);
                fmt.Println("salary = ", element.Salary);
        }
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


