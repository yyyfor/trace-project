package factory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

type FactoryChaincode struct {
}

type factoryInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (t *FactoryChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Factory success init."))
}

func (t *FactoryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "putFactoryInfo" {
		return t.putFactory(stub, args)
	} else if funcName == "getFactoryInfo" {
		return t.getFactory(stub, args)
	} else if funcName == "getFactoryHistory" {
		return t.getFactoryHistory(stub, args)
	}

	return shim.Error(fmt.Sprintf("no such method."))
}

func (t *FactoryChaincode) putFactory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Producer parameter should be 3.")
	}

	id := args[0]
	name := args[1]

	factoryAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprint(err.Error()))
	}

	info := factoryInfo{}
	if factoryAsJSON != nil {
		err = json.Unmarshal(factoryAsJSON, &info)
		if err != nil {
			shim.Error(err.Error())
		}
		info.Name = name
	} else {
		i, err := strconv.Atoi(id)
		if err != nil {
			shim.Error(err.Error())
		} else {
			info = factoryInfo{ID: i, Name: name}
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

	return shim.Success(infoAsJSONBytes)

}

func (t *FactoryChaincode) getFactory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("parameter should be exact 1.")
	}
	id := args[0]
	infoAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if infoAsJSON == nil {
		return shim.Error(fmt.Sprintf("id is not exist"))
	}
	return shim.Success(infoAsJSON)
}

func (t *FactoryChaincode) getFactoryHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	err := shim.Start(new(FactoryChaincode))
	if err != nil {
		fmt.Printf("Error starting producer chaincode")
	}
}
