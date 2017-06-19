/* 
    chaincode for sharing history storage

    author: Runshan Hu
*/

package main

import (
        "errors"
        "fmt"
        "math"
        "log"
        "encoding/json"
        "net/http"
        "io/ioutil"
        "github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

// value format for ledger
type ledgerMes struct {
  RemainBudget     float64   `json:"budget"`
  FunType          []string  `json:"funType"`
  Result           []float64 `json:"results"`
}

// message format for query
type queryMes struct {
  RequestBudget    float64   `json:"budget"`
  FunType          string    `json:"funType"`
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
        } else if function == "query" {
                return t.queryMatchTest(stub, args)
        }

        fmt.Println("invoke did not find func:" + function)

        return nil, errors.New("Received unknown function invocation")

}


func (t *SimpleChaincode) queryMatchTest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
       var dataId string
       var value string
       var err error

       utility_bound := 0.1

       // args should have two parameter: datasetId and user's query
       if len(args) != 2 {
               return nil, errors.New("Incorrect number of arguments. Expecting 2. DatasetID and your query ")
       }
       
       dataId = args[0];
       value = args[1];
       
       //parser user's query 
       mes_from_query := queryMes{}
       json.Unmarshal([]byte(value), &mes_from_query);

       // get the old query from ledger
       valAsbytes, err := stub.GetState(dataId);
       if err != nil {
                jsonResp := "{\"Error\": \"Failed to get the state for " + dataId + "\"}"
                return nil, errors.New(jsonResp) 
       }

       // parser the old query (from ledger)
       mes_from_ledger := ledgerMes{}
       json.Unmarshal(valAsbytes, &mes_from_ledger)

       flag := false;   //whether query exists before
       var old_result, final_result, perturbed_result float64
       var i int
       var e string
       smallbudget := 0.1;
       for i, e = range mes_from_ledger.FunType {
                if e == mes_from_query.FunType {
                        flag = true 
                        break 
                } 
       }
       // if old result exist
       if flag {
                // old result (from ledger)
                old_result = mes_from_ledger.Result[i]
                // get perturbed result from anonymisation service
                perturbed_result = getResultAnonyService(mes_from_query.FunType, smallbudget)
                // utility test
                if math.Abs(old_result - perturbed_result) < utility_bound {
                        final_result =  perturbed_result
                        updateLedger(stub, dataId, mes_from_query.FunType, final_result, smallbudget)
                        
                } else {
                       if mes_from_ledger.RemainBudget >= mes_from_query.RequestBudget  {
                        final_result = getResultAnonyService(mes_from_query.FunType, mes_from_query.RequestBudget)
                        // updateLedger()
                        updateLedger(stub, dataId, mes_from_query.FunType, final_result, mes_from_query.RequestBudget)
                      } else {
                        final_result = -1000 
                        // updateLedger()
                        updateLedger(stub, dataId, mes_from_query.FunType, final_result, smallbudget)
                      }
                }
       } else { // old result not exist
                       if mes_from_ledger.RemainBudget >= mes_from_query.RequestBudget  {
                               final_result = getResultAnonyService(mes_from_query.FunType, mes_from_query.RequestBudget)
                               //updateLedger()
                               updateLedger(stub, dataId, mes_from_query.FunType, final_result, mes_from_query.RequestBudget)
                       } else {
                               final_result = -1000 
                               // updateLedger()
                               updateLedger(stub, dataId, mes_from_query.FunType, final_result, 0)
                      }
       }
       return nil, nil
}

func updateLedger(stub shim.ChaincodeStubInterface, dataId string, funType string, newResult float64, subBudget float64) ([]byte, error) {

        valAsbytes, err := stub.GetState(dataId)
        if err != nil {
          jsonResp := "{\"Error\": \"Failed to get the state for " + dataId + "\"}"
                 return nil, errors.New(jsonResp) 
        }
        
        newValue := ledgerMes{} 
        json.Unmarshal(valAsbytes, &newValue)
        newValue.RemainBudget = newValue.RemainBudget - subBudget

        var index int
        for i,e := range newValue.FunType {
                if e == funType {
                        index = i
                        break
                }
        }
        newValue.Result[index] = newResult
        newValue_json,err := json.Marshal(newValue)
        
        // write back to the ledger
        err = stub.PutState(dataId, []byte(newValue_json))
        if err != nil {
               return nil, err 
        }
        return nil, nil
}

type serviceResult struct {
         Result float64 `json:"result"`
}

func getResultAnonyService( funtype string, budget float64  ) float64  {
 
        resp, err := http.Get("http://10.7.6.25:3000/dataset/sum")
        switch funtype {
               case "sum": 
                         resp, err = http.Get("http://10.7.6.25:3000/dataset/sum")
               case "avg": 
                         resp, err = http.Get("http://10.7.6.25:3000/dataset/avg")
               case "max": 
                         resp, err = http.Get("http://10.7.6.25:3000/dataset/max")
               case "min": 
                         resp, err = http.Get("http://10.7.6.25:3000/dataset/min")
        /*       default:
                       {
                         log.Println("unrecognized function type")
                         return nil
                       }
       */
        } 
        if err != nil {
                log.Println(err);
        }

        body,err := ioutil.ReadAll(resp.Body);
        if err != nil {
                log.Println(err);
        }

        defer resp.Body.Close();

        result := serviceResult{}

        if err := json.Unmarshal(body, &result); err != nil {
                log.Println(err); 
        }

        return result.Result;
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


