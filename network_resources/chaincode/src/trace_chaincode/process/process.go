package process

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

type producerInfo struct {
	ID      int    `json:"id"`
	Product string `json:"product"`
	Source  string `json:"source"`
}

type factoryInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	PrevID int    `json:"prevID"`
}

type saleInfo struct {
	ID     int  `json:"id"`
	PrevID int  `json:"prevID"`
	IsSale bool `json:"isSale"`
}

func (t *Processor) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Processor success init."))
}

func (t *Processor) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "putProducerInfo" {
		return t.putProducerInfo(stub, args)
	} else if funcName == "getProducerInfo" {
		return t.getProducerInfo(stub, args)
	} else if funcName == "getProducerHistory" {
		return t.getProducerHistory(stub, args)
	} else if funcName == "putFactoryInfo" {
		return t.putFactoryInfo(stub, args)
	} else if funcName == "getFactoryInfo" {
		return t.getProducerInfo(stub, args)
	} else if funcName == "getFactoryHistory" {
		return t.getFactoryHistory(stub, args)
	}

	return shim.Error(fmt.Sprintf("no such method."))
}

func (t *Processor) putProducerInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Producer parameter should be 3.")
	}

	id := args[0]
	productName := args[1]
	source := args[2]

	producerAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprint(err.Error()))
	}

	info := producerInfo{}
	if producerAsJSON != nil {
		err = json.Unmarshal(producerAsJSON, &info)
		if err != nil {
			shim.Error(err.Error())
		}
		info.Product = productName
		info.Source = source
	} else {
		i, err := strconv.Atoi(id)
		if err != nil {
			shim.Error(err.Error())
		} else {
			info = producerInfo{ID: i, Product: productName, Source: source}
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

	return shim.Success([]byte("set producer info successful"))

}

func (t *Processor) getProducerInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("parameter should be exact 1.")
	}
	id := args[0]
	infoAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if infoAsJSON == nil {
		return shim.Error(fmt.Sprintf("procedure id is not exist"))
	}
	return shim.Success(infoAsJSON)
}

func (t *Processor) getProducerHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		tm := time.Unix(txtimestamp.Seconds, int64(txtimestamp.Nanos)).String()
		fmt.Printf("Tx info - txid:%s value: %s if delete:%t\ndatetime: %s", txid, string(txvalue), txstatus, tm)
		keys = append(keys, string(txvalue)+":"+tm)
	}

	jsonIds, err := json.Marshal(keys)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonIds)
}

func (t *Processor) putFactoryInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Producer parameter should be 3.")
	}

	id := args[0]
	name := args[1]
	prevID := args[2]
	prevIDNum, err := strconv.Atoi(prevID)
	if err != nil {
		return shim.Error(err.Error())
	}

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
		info.PrevID = prevIDNum
	} else {
		i, err := strconv.Atoi(id)
		if err != nil {
			shim.Error(err.Error())
		} else {
			info = factoryInfo{ID: i, Name: name, PrevID: prevIDNum}
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

	return shim.Success([]byte("set factory info successful"))

}

func (t *Processor) getFactoryInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("parameter should be exact 1.")
	}
	id := args[0]
	infoAsJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	} else if infoAsJSON == nil {
		return shim.Error(fmt.Sprintf("factory id is not exist"))
	}
	factory := factoryInfo{}
	err = json.Unmarshal(infoAsJSON, &factory)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(infoAsJSON)
}

func (t *Processor) getFactoryHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("parameter should be exact 1.")
	}
	id := args[0]

	infoAsJSON, err := stub.GetState(id)
	factory := factoryInfo{}
	var procedureHistory string

	if infoAsJSON != nil {
		err := json.Unmarshal(infoAsJSON, &factory)
		if err != nil {
			return shim.Error(err.Error())
		}
		producerID := factory.ID
		producerID_str := strconv.Itoa(producerID)
		producerInfo := t.getProducerHistory(stub, []string{producerID_str})
		procedureHistory = producerInfo.GetMessage()

	} else {
		return shim.Error("factory is not exist.")
	}

	idIter, err := stub.GetHistoryForKey(id)
	if err != nil {
		shim.Error(err.Error())
	}

	defer idIter.Close()

	if idIter == nil {
		return shim.Error("history is not exists.")
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
		tm := time.Unix(txtimestamp.Seconds, int64(txtimestamp.Nanos)).String()
		fmt.Printf("Tx info - txid:%s value: %s if delete:%t\ndatetime: %s", txid, string(txvalue), txstatus, tm)
		keys = append(keys, string(txvalue)+":"+tm)
	}
	keys = append(keys, procedureHistory)

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
