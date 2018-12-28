package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

type TraceChaincode struct {
}

type Product struct {
	ID          int    `json:"id"`
	ProductName string `json:"productName"`
	Factory     string `json:"factory"`
}

func (t *TraceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("producer success init."))
}

func (t *TraceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "putProduct" {
		return t.putProduct(stub, args)
	} else if funcName == "getProduct" {
		return t.getProduct(stub, args)
	} else if funcName == "getProductHistory" {
		return t.getProductHistory(stub, args)
	}

	return shim.Error(fmt.Sprintf("no such method."))
}

func (t *TraceChaincode) putProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	productName := args[1]
	factory := args[2]

	productAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprint(err.Error()))
	}

	product := Product{}
	if productAsJSON != nil {
		err = json.Unmarshal(productAsJSON, &product)
		if err != nil {
			shim.Error(err.Error())
		}
		product.ProductName = productName
		product.Factory = factory
	} else {
		i, err := strconv.Atoi(id)
		if err != nil {
			shim.Error(err.Error())
		} else {
			product = Product{ID: i, ProductName: productName, Factory: factory}
		}
	}

	productAsJSONBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(id, productAsJSONBytes)
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success([]byte(productAsJSONBytes))

}

func (t *TraceChaincode) getProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	productAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if productAsJSON == nil {
		return shim.Error(fmt.Sprintf("product is not exist"))
	}
	return shim.Success(productAsJSON)
}

func (t *TraceChaincode) getProductHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("parameter should be exact 1.")
	}
	id := args[0]
	idIter, err := stub.GetHistoryForKey(id)
	if err != nil {
		shim.Error(err.Error())
	}

	defer idIter.Close()

	if idIter == nil {
		return shim.Error("history is not exists.")
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for idIter.HasNext() {
		response, err := idIter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	if !bArrayMemberAlreadyWritten {
		buffer.WriteString("No record found")
	}
	buffer.WriteString("]")

	fmt.Printf("- getAllTransactionForNumber returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(TraceChaincode))
	if err != nil {
		fmt.Printf("Error starting producer chaincode")
	}
}
