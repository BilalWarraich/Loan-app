/*
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//SmartContract is the data structure which represents this contract and on which  various contract lifecycle functions are attached
type SmartContract struct {
}

type Lender struct {
	ObjectType string `json:"Type"`
	LenderID   string `json:"lenderID"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type Lendee struct {
	ObjectType string `json:"Type"`
	LendeeID   string `json:"lendeeID"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type Request struct {
	ObjectType      string `json:"Type"`
	RequestID       string `json:"requestID"`
	LendeeID        string `json:"lendeeID"`
	LenderID        string `json:"lenderID"`
	LoanRequired    string `json:"loanRequired"`
	Product         string `json:"product"`
	ReturnDate      string `json:"returnDate"`
	ExpectedOutcome string `json:"expectedOutcome"`
	Status          string `json:"status"`
}

type Response struct {
	ObjectType string `json:"Type"`
	ResponseID string `json:"responseID"`
	LendeeID   string `json:"lendeeID"`
	LenderID   string `json:"lenderID"`
	LoanAmount string `json:"loanAmount"`
	TimeLimit  string `json:"timeLimit"`
	Interest   string `json:"interest"`
	Status     string `json:"status"`
}

func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {

	fmt.Println("Init Firing!")
	return shim.Success(nil)
}

func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Chaincode Invoke Is Running " + function)
	if function == "addLender" {
		return t.addLender(stub, args)
	}
	if function == "queryLender" {
		return t.queryLender(stub, args)
	}
	if function == "queryAllLenders" {
		return t.queryAllLenders(stub, args)
	}
	if function == "queryLenderByName" {
		return t.queryLenderByName(stub, args)
	}
	if function == "queryLenderByID" {
		return t.queryLenderByID(stub, args)
	}
	if function == "addLendee" {
		return t.addLendee(stub, args)
	}
	if function == "queryLendee" {
		return t.queryLendee(stub, args)
	}
	if function == "queryLendeeByID" {
		return t.queryLendeeByID(stub, args)
	}
	if function == "addRequest" {
		return t.addRequest(stub, args)
	}
	if function == "queryRequest" {
		return t.queryRequest(stub, args)
	}
	if function == "queryRequests" {
		return t.queryRequests(stub, args)
	}
	if function == "queryRequestByID" {
		return t.queryRequestByID(stub, args)
	}
	if function == "queryRequestID" {
		return t.queryRequestID(stub, args)
	}
	if function == "addResponse" {
		return t.addResponse(stub, args)
	}
	if function == "queryResponse" {
		return t.queryResponse(stub, args)
	}
	if function == "queryResponseByID" {
		return t.queryResponseByID(stub, args)
	}
	if function == "queryResponseID" {
		return t.queryResponseID(stub, args)
	}
	if function == "updateRequest" {
		return t.updateRequest(stub, args)
	}
	if function == "updateResponse" {
		return t.updateResponse(stub, args)
	}

	fmt.Println("Invoke did not find specified function " + function)
	return shim.Error("Invoke did not find specified function " + function)
}

func (t *SmartContract) addLender(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect Number of Aruments. Expecting 3")
	}

	fmt.Println("Adding new Lender")

	// ==== Input sanitation ====
	if len(args[0]) <= 0 {
		return shim.Error("1st Argument Must be a Non-Empty String")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd Argument Must be a Non-Empty String")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd Argument Must be a Non-Empty String")
	}

	lenderID := args[0]
	username := args[1]
	password := args[2]

	// ======Check if lender Already exists

	lenderAsBytes, err := stub.GetState(lenderID)
	if err != nil {
		return shim.Error("Transaction Failed with Error: " + err.Error())
	} else if lenderAsBytes != nil {
		return shim.Error("The Inserted Lender ID already Exists: " + lenderID)
	}

	// ======Check if lender Already exists

	lenderAsName, err := stub.GetState(username)
	if err != nil {
		return shim.Error("Transaction Failed with Error: " + err.Error())
	} else if lenderAsName != nil {
		return shim.Error("The Inserted Lender username already Exists: " + username)
	}

	// ===== Create lender Object and Marshal to JSON

	objectType := "lender"
	lender := &Lender{objectType, lenderID, username, password}
	lenderJSONasBytes, err := json.Marshal(lender)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Save lender to State

	err = stub.PutState(lenderID, lenderJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Return Success

	fmt.Println("Successfully Saved lender")
	return shim.Success(nil)
}

func (t *SmartContract) queryLender(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	username := args[0]
	password := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"lender\",\"username\":\"%s\",\"password\":\"%s\"}}", username, password)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryLenderByName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	username := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"lender\",\"username\":\"%s\"}}", username)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryLenderByID(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lenderID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"lender\",\"lenderID\":\"%s\"}}", lenderID)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryAllLenders(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"lender\"}}")

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) addLendee(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect Number of Aruments. Expecting 3")
	}

	fmt.Println("Adding new Lendee")

	// ==== Input sanitation ====
	if len(args[0]) <= 0 {
		return shim.Error("1st Argument Must be a Non-Empty String")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd Argument Must be a Non-Empty String")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd Argument Must be a Non-Empty String")
	}

	lendeeID := args[0]
	username := args[1]
	password := args[2]

	// ======Check if lendee Already exists

	lendeeAsBytes, err := stub.GetState(lendeeID)
	if err != nil {
		return shim.Error("Transaction Failed with Error: " + err.Error())
	} else if lendeeAsBytes != nil {
		return shim.Error("The Inserted Lendee ID already Exists: " + lendeeID)
	}

	// ===== Create lendee Object and Marshal to JSON

	objectType := "lendee"
	lendee := &Lendee{objectType, lendeeID, username, password}
	lendeeJSONasBytes, err := json.Marshal(lendee)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Save lendee to State

	err = stub.PutState(lendeeID, lendeeJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Return Success

	fmt.Println("Successfully Saved lendee")
	return shim.Success(nil)
}

func (t *SmartContract) queryLendee(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	username := args[0]
	password := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"lendee\",\"username\":\"%s\",\"password\":\"%s\"}}", username, password)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryLendeeByID(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lendeeID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"lendee\",\"lendeeID\":\"%s\"}}", lendeeID)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) addRequest(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var err error

	if len(args) != 8 {
		return shim.Error("Incorrect Number of Aruments. Expecting 8")
	}

	fmt.Println("Adding new Request")

	// ==== Input sanitation ====
	if len(args[0]) <= 0 {
		return shim.Error("1st Argument Must be a Non-Empty String")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd Argument Must be a Non-Empty String")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd Argument Must be a Non-Empty String")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th Argument Must be a Non-Empty String")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th Argument Must be a Non-Empty String")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th Argument Must be a Non-Empty String")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th Argument Must be a Non-Empty String")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8th Argument Must be a Non-Empty String")
	}

	requestID := args[0]
	lendeeID := args[1]
	lenderID := args[2]
	loanRequired := args[3]
	product := args[4]
	returnDate := args[5]
	expectedOutcome := args[6]
	status := args[7]

	// ======Check if request Already exists

	requestAsBytes, err := stub.GetState(requestID)
	if err != nil {
		return shim.Error("Transaction Failed with Error: " + err.Error())
	} else if requestAsBytes != nil {
		return shim.Error("The Inserted requestID ID already Exists: " + requestID)
	}

	// ===== Create request Object and Marshal to JSON

	objectType := "request"
	request := &Request{objectType, requestID, lendeeID, lenderID, loanRequired, product, returnDate, expectedOutcome, status}
	requestJSONasBytes, err := json.Marshal(request)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Save request to State

	err = stub.PutState(requestID, requestJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Return Success

	fmt.Println("Successfully Saved request")
	return shim.Success(nil)
}

func (t *SmartContract) queryRequestByID(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lendeeID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"request\",\"lendeeID\":\"%s\"}}", lendeeID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryRequest(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	status := args[0]
	lenderID := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"request\",\"status\":\"%s\",\"lenderID\":\"%s\"}}", status, lenderID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryRequests(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lendeeID := args[0]
	lenderID := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"request\",\"lendeeID\":\"%s\",\"lenderID\":\"%s\"}}", lendeeID, lenderID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryRequestID(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	requestID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"request\",\"requestID\":\"%s\"}}", requestID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) addResponse(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var err error

	if len(args) != 7 {
		return shim.Error("Incorrect Number of Aruments. Expecting 7")
	}

	fmt.Println("Adding new Response")

	// ==== Input sanitation ====
	if len(args[0]) <= 0 {
		return shim.Error("1st Argument Must be a Non-Empty String")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd Argument Must be a Non-Empty String")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd Argument Must be a Non-Empty String")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th Argument Must be a Non-Empty String")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th Argument Must be a Non-Empty String")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th Argument Must be a Non-Empty String")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th Argument Must be a Non-Empty String")
	}

	responseID := args[0]
	lendeeID := args[1]
	lenderID := args[2]
	loanAmount := args[3]
	timeLimit := args[4]
	interest := args[5]
	status := args[6]

	// ======Check if response Already exists

	responseAsBytes, err := stub.GetState(responseID)
	if err != nil {
		return shim.Error("Transaction Failed with Error: " + err.Error())
	} else if responseAsBytes != nil {
		return shim.Error("The Inserted responseID ID already Exists: " + responseID)
	}

	// ===== Create response Object and Marshal to JSON

	objectType := "response"
	response := &Response{objectType, responseID, lendeeID, lenderID, loanAmount, timeLimit, interest, status}
	responseJSONasBytes, err := json.Marshal(response)

	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Save response to State

	err = stub.PutState(responseID, responseJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ======= Return Success

	fmt.Println("Successfully Saved response")
	return shim.Success(nil)
}

func (t *SmartContract) queryResponse(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	status := args[0]
	lendeeID := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"response\",\"status\":\"%s\",\"lendeeID\":\"%s\"}}", status, lendeeID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryResponseByID(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lenderID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"response\",\"lenderID\":\"%s\"}}", lenderID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) queryResponseID(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	responseID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"response\",\"responseID\":\"%s\"}}", responseID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SmartContract) updateRequest(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	requestID := args[0]
	newStatus := args[1]
	fmt.Println("- start  ", requestID, newStatus)

	requestAsBytes, err := stub.GetState(requestID)
	if err != nil {
		return shim.Error("Failed to get status:" + err.Error())
	} else if requestAsBytes == nil {
		return shim.Error("Request does not exist")
	}

	requestToUpdate := Request{}
	err = json.Unmarshal(requestAsBytes, &requestToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	requestToUpdate.Status = newStatus //change the status

	requestJSONasBytes, _ := json.Marshal(requestToUpdate)
	err = stub.PutState(requestID, requestJSONasBytes) //rewrite
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end  (success)")
	return shim.Success(nil)
}

func (t *SmartContract) updateResponse(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	responseID := args[0]
	newStatus := args[1]
	fmt.Println("- start  ", responseID, newStatus)

	responseAsBytes, err := stub.GetState(responseID)
	if err != nil {
		return shim.Error("Failed to get status:" + err.Error())
	} else if responseAsBytes == nil {
		return shim.Error("response does not exist")
	}

	responseToUpdate := Response{}
	err = json.Unmarshal(responseAsBytes, &responseToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	responseToUpdate.Status = newStatus //change the status

	responseJSONasBytes, _ := json.Marshal(responseToUpdate)
	err = stub.PutState(responseID, responseJSONasBytes) //rewrite
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end  (success)")
	return shim.Success(nil)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

//Main Function starts up the Chaincode
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Smart Contract could not be run. Error Occured: %s", err)
	} else {
		fmt.Println("Smart Contract successfully Initiated")
	}
}
