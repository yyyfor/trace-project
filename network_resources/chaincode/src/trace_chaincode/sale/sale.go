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

type SaleChaincode struct {
}

type Info struct {
	ID     int    `json:"id"`
	Store  string `json:"store"`
	IsSale bool   `json:"isSale"`
}

func (t *SaleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Processor success init."))
}

func (t *SaleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "putSaleInfo" {
		return t.putInfo(stub, args)
	} else if funcName == "getSaleInfo" {
		return t.getInfo(stub, args)
	} else if funcName == "getIdHistory" {
		return t.getIdHistory(stub, args)
	}

	return shim.Error(fmt.Sprintf("no such method."))
}

func (t *SaleChaincode) putInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	storeName := args[1]
	isSale := args[2]

	b, err := strconv.ParseBool(isSale)
	if err != nil {
		shim.Error(err.Error())
	}
	infoAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprint(err.Error()))
	}

	info := Info{}
	if infoAsJSON != nil {
		err = json.Unmarshal(infoAsJSON, &info)
		if err != nil {
			shim.Error(err.Error())
		}
		info.Store = storeName
		info.IsSale = b
	} else {
		i, err := strconv.Atoi(id)
		if err != nil {
			shim.Error(err.Error())
		} else {
			info = Info{ID: i, Store: storeName, IsSale: b}
		}
	}

	infoAsJSONBytes, err := json.Marshal(info)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(id, infoAsJSONBytes)
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success([]byte(infoAsJSONBytes))

}

func (t *SaleChaincode) getInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	infoAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if infoAsJSON == nil {
		return shim.Error(fmt.Sprintf("sale record is not exist"))
	}
	return shim.Success(infoAsJSON)
}

func (t *SaleChaincode) getIdHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

// func (t *Sale) getIdAllHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	id := args[0]
// 	invokeArgs := util.ToChaincodeArgs("invoke", "getIdHistory", id)
// 	response := stub.InvokeChaincode("producer", invokeArgs, "trace_channel")
// 	if response.Status != shim.OK {
// 		errMsg := fmt.Sprint("Failed to invoke chaincode. error is %s", response.Payload)
// 		return shim.Error(errMsg)
// 	}

// 	return shim.Success(response.Payload)

// }

func main() {
	err := shim.Start(new(SaleChaincode))
	if err != nil {
		fmt.Printf("Error starting producer chaincode")
	}
}
