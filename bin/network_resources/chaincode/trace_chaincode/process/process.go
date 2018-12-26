package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

type Processor struct {
}

type Info struct {
	ID      int    `json:"id"`
	Product string `json:"product"`
	Factory string `json:"factory"`
}

func (t *Processor) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Processor success init."))
}

func (t *Processor) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "putInfo" {
		return t.putInfo(stub, args)
	} else if funcName == "getInfo" {
		return t.getInfo(stub, args)
	} else if funcName == "getIdHistory" {
		return t.getIdHistory(stub, args)
	}

	return shim.Error(fmt.Sprintf("no such method."))
}

func (t *Processor) putInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	productName := args[1]
	factory := args[2]

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
		info.Product = productName
		info.Factory = factory
	} else {
		i, err := strconv.Atoi(id)
		if err != nil {
			shim.Error(err.Error())
		} else {
			info = Info{ID: i, Product: productName, Factory: factory}
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

	return shim.Success([]byte("put info successful"))

}

func (t *Processor) getInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	infoAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if infoAsJSON == nil {
		return shim.Error(fmt.Sprintf("process info is not exist"))
	}
	return shim.Success(infoAsJSON)
}

func (t *Processor) getIdHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	idIter, err := stub.GetHistoryForKey(id)
	if err != nil {
		shim.Error(err.Error())
	}

	defer idIter.Close()

	if idIter == nil {
		return shim.Success([]byte("no history."))
	}

	var keys []string
	for idIter.HasNext() {
		response, err := idIter.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("Get History For id operaion failed. Error state is: %s", err))
		}
		txid := response.TxId
		txvalue := response.Value
		txstatus := response.IsDelete
		txtimestamp := response.Timestamp
		tm := time.Unix(txtimestamp.Seconds, 0)
		dateStr := tm.Format("2000-01-01 01:01:01 AM")
		fmt.Printf("Tx info - txid:%s value: %s if delete:%t\ndatetime: %s", txid, string(txvalue), txstatus, dateStr)
		keys = append(keys, string(txvalue)+":"+dateStr)
	}

	jsonIds, err := json.Marshal(keys)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonIds)
}

func main() {
	err := shim.Start(new(Processor))
	if err != nil {
		fmt.Printf("Error starting producer chaincode")
	}
}
