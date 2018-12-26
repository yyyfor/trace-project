package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

type Sale struct {
}

type Info struct {
	ID     int    `json:"id"`
	Store  string `json:"store"`
	IsSale bool   `json:"isSale"`
}

func (t *Sale) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Processor success init."))
}

func (t *Sale) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "putInfo" {
		return t.putInfo(stub, args)
	} else if funcName == "getInfo" {
		return t.getInfo(stub, args)
	} else if funcName == "getIdHistory" {
		return t.getIdHistory(stub, args)
	} else if funcName == "getIdAllHistory" {
		return t.getIdAllHistory(stub, args)
	}

	return shim.Error(fmt.Sprintf("no such method."))
}

func (t *Sale) putInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	return shim.Success([]byte("put info successful"))

}

func (t *Sale) getInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	infoAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if infoAsJSON == nil {
		return shim.Error(fmt.Sprintf("sale record is not exist"))
	}
	return shim.Success(infoAsJSON)
}

func (t *Sale) getIdHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

func (t *Sale) getIdAllHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	invokeArgs := util.ToChaincodeArgs("invoke", "getIdHistory", id)
	response := stub.InvokeChaincode("producer", invokeArgs, "trace_channel")
	if response.Status != shim.OK {
		errMsg := fmt.Sprint("Failed to invoke chaincode. error is %s", response.Payload)
		return shim.Error(errMsg)
	}

	return shim.Success(response.Payload)

}

func main() {
	err := shim.Start(new(Sale))
	if err != nil {
		fmt.Printf("Error starting producer chaincode")
	}
}
