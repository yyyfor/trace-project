package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

type Producer struct {
}

type Product struct {
	ID          int    `json:"id"`
	ProductName string `json:"productName"`
	Factory     string `json:"factory"`
}

func (t *Producer) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("producer success init."))
}

func (t *Producer) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "putProduct" {
		return t.putProduct(stub, args)
	} else if funcName == "getProduct" {
		return t.getProduct(stub, args)
	} else if funcName == "getIdHistory" {
		return t.getIdHistory(stub, args)
	}

	return shim.Error(fmt.Sprintf("no such method."))
}

func (t *Producer) putProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	return shim.Success([]byte("put product successful"))

}

func (t *Producer) getProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id := args[0]
	productAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if productAsJSON == nil {
		return shim.Error(fmt.Sprintf("product is not exist"))
	}
	return shim.Success(productAsJSON)
}

func (t *Producer) getIdHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	err := shim.Start(new(Producer))
	if err != nil {
		fmt.Printf("Error starting producer chaincode")
	}
}
