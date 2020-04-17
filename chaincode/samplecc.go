package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SampleChaincode example sample Chaincode implementation
type SampleChaincode struct {
}

type car struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Owner      string `json:"owner"`    //the fieldtags are needed to keep case from bouncing around
	Color      string `json:"color"`
	Model      string `json:"model"`
	
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Printf("Error starting Sample chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initCar" { //create a new car
		return t.initCar(stub, args)
	} else if function == "transferCar" { //change owner of a specific car
		return t.transferCar(stub, args)
	} else if function == "delete" { //delete a car
		return t.delete(stub, args)
	} else if function == "readCar" { //read a car
		return t.readCar(stub, args) 	
	} else if function =="initLedger" {
		return t.initLedger(stub)
	}else if function == "getAllCars"{
		return t.getAllCars(stub)
	}else if function == "getCarsByRange" {
		return t.getCarsByRange(stub,args)
	}



	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initCar - create a new car, store into chaincode state
// ============================================================
func (t *SampleChaincode) initCar(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1       2     
	// "Ashwani", "Blue", "BMW"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init car")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	
	owner := args[0]
	color := strings.ToLower(args[1])
	model := strings.ToLower(args[2])
	

	// ==== Check if car already exists ====
	carAsBytes, err := stub.GetState(owner)
	if err != nil {
		return shim.Error("Failed to get car: " + err.Error())
	} else if carAsBytes != nil {
		fmt.Println("This car already exists for owner: " + owner)
		return shim.Error("This car already exists: " + owner)
	}

	// ==== Create car object and marshal to JSON ====
	objectType := "Car"
	car := &car{objectType, owner, color, model}
	carJSONasBytes, err := json.Marshal(car)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the car json string manually if you don't want to use struct marshalling
	//carJSONasString := `{"docType":"Car",  "owner": "` + owner + `", "color": "` + color + `", "model": ` + model + `"}`
	//carJSONasBytes := []byte(carJSONasString)

	// === Save car to state ===
	err = stub.PutState(owner, carJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	// ==== Car saved. Return success ====
	fmt.Println("- end init car")
	return shim.Success(nil)
}

// ===============================================
// initLedger - populate the ledger with Cars
// ===============================================
func (s *SampleChaincode) initLedger(stub shim.ChaincodeStubInterface) pb.Response {
        cars := []car{
                car{Owner: "Ashwani", Color: "blue",Model: "bmw" },
		car{Owner: "Raja", Color: "red",Model: "santro" },
		car{Owner: "Naman", Color: "white",Model: "wagonR" },
		car{Owner: "Gurneet", Color: "black",Model: "fortuner" },
		car{Owner: "John", Color: "orange",Model: "pajero" },
		car{Owner: "Bob", Color: "grey",Model: "RangeRover" },
		car{Owner: "Alice", Color: "voilet",Model: "bentley" },
		car{Owner: "Mayank", Color: "ruby",Model: "merc" },
		car{Owner: "Ratan", Color: "green",Model: "audi" },
		car{Owner: "Prem", Color: "seagreen",Model: "kia" },
		car{Owner: "Yoyo", Color: "silver",Model: "thar" },
		car{Owner: "Mini", Color: "gold",Model: "xuv" },
				
                
        }

        i := 0
        for i < len(cars) {
                fmt.Println("i is ", i)
                carAsBytes, _ := json.Marshal(cars[i])
                stub.PutState("CAR"+strconv.Itoa(i), carAsBytes)
                fmt.Println("Added", cars[i])
                i = i + 1
        }

        return shim.Success(nil)
}

// ===============================================
// getAllCars - populate the ledger with Cars
// ===============================================


func (s *SampleChaincode) getAllCars(APIstub shim.ChaincodeStubInterface) pb.Response {

        startKey := ""
        endKey := ""

        resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()

        // buffer is a JSON array containing QueryResults
        var buffer bytes.Buffer
        buffer.WriteString("[")

        bArrayMemberAlreadyWritten := false
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                // Add a comma before array members, suppress it for the first array member
                if bArrayMemberAlreadyWritten == true {
                        buffer.WriteString(",")
                }
                buffer.WriteString("{\"Key\":")
                buffer.WriteString("\"")
                buffer.WriteString(queryResponse.Key)
                buffer.WriteString("\"")

                buffer.WriteString(", \"Record\":")
                // Record is a JSON object, so we write as-is
                buffer.WriteString(string(queryResponse.Value))
                buffer.WriteString("}")
  }
        buffer.WriteString("]")

        fmt.Printf("- getAllCars:\n%s\n", buffer.String())

        return shim.Success(buffer.Bytes())
}


// ===============================================
// getCarsByRange- get Cars by Range
// ===============================================


func (s *SampleChaincode) getCarsByRange(APIstub shim.ChaincodeStubInterface, args[] string) pb.Response {


	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// ==== Input sanitation ====
	fmt.Println("- start getCarsByRange")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	
	
        startKey := args[0]
        endKey := args[1]

        resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()

        // buffer is a JSON array containing QueryResults
        var buffer bytes.Buffer
        buffer.WriteString("[")

        bArrayMemberAlreadyWritten := false
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                // Add a comma before array members, suppress it for the first array member
                if bArrayMemberAlreadyWritten == true {
                        buffer.WriteString(",")
                }
                buffer.WriteString("{\"Key\":")
                buffer.WriteString("\"")
                buffer.WriteString(queryResponse.Key)
                buffer.WriteString("\"")

                buffer.WriteString(", \"Record\":")
                // Record is a JSON object, so we write as-is
                buffer.WriteString(string(queryResponse.Value))
                buffer.WriteString("}")
  }
        buffer.WriteString("]")

        fmt.Printf("- getAllCars:\n%s\n", buffer.String())

        return shim.Success(buffer.Bytes())
}



// ===============================================
// readCar - read a car from chaincode state
// ===============================================
func (t *SampleChaincode) readCar(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the owner to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the car from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Car does not exist for owner: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a car key/value pair from state
// ==================================================
func (t *SampleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var carJSON car
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	owner := args[0]

	
	valAsbytes, err := stub.GetState(owner) //get the ownerName from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + owner + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Car does not exist for owner: " + owner + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &carJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of car: " + owner + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(owner) //remove the car from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}


	return shim.Success(nil)
}

// ===========================================================
// transfer a car by setting a new owner name on the car
// ===========================================================
func (t *SampleChaincode) transferCar(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	owner := args[0]
	newOwner := strings.ToLower(args[1])
	fmt.Println("- start transferCar ", owner, newOwner)

	carAsBytes, err := stub.GetState(owner)
	if err != nil {
		return shim.Error("Failed to get car:" + err.Error())
	} else if carAsBytes == nil {
		return shim.Error("Car does not exist")
	}

	carToTransfer := car{}
	err = json.Unmarshal(carAsBytes, &carToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	carToTransfer.Owner = newOwner //change the owner

	carJSONasBytes, _ := json.Marshal(carToTransfer)
	err = stub.PutState(owner, carJSONasBytes) //rewrite the car
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferCar (success)")
	return shim.Success(nil)
}


