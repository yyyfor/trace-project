package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkGetProducerInfo(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getProducerInfo"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("getInfo", id, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("getInfo", id, "failed to get value")
		t.FailNow()
	}
	fmt.Println(string(res.Payload))
}

func checkPutProducerInfo(t *testing.T, stub *shim.MockStub, id string, productName string, factory string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("putProducerInfo"), []byte(id), []byte(productName), []byte(factory)})
	if res.Status != shim.OK {
		fmt.Println("putInfo", id, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkGetProducerHistory(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getProcedureHistory"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("getIdHistory", id, "failed", string(res.Message))
		t.FailNow()
	}
	fmt.Println(string(res.Payload))
}

func checkGetFactoryInfo(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getFactoryInfo"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("getInfo", id, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("getInfo", id, "failed to get value")
		t.FailNow()
	}
	fmt.Println(string(res.Payload))
}

func checkPutFactoryInfo(t *testing.T, stub *shim.MockStub, id string, productName string, factory string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("putFactoryInfo"), []byte(id), []byte(productName), []byte(factory)})
	if res.Status != shim.OK {
		fmt.Println("putInfo", id, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkGetFactoryHistory(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getFactoryHistory"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("getIdHistory", id, "failed", string(res.Message))
		t.FailNow()
	}
	fmt.Println(string(res.Payload))
}

func Test_Invoke(t *testing.T) {
	scc := new(Processor)
	stub := shim.NewMockStub("test1", scc)

	checkInit(t, stub, [][]byte{[]byte("test")})

	// Invoke A->B for 123
	checkPutProducerInfo(t, stub, "1", "pig", "fusheng")
	checkPutProducerInfo(t, stub, "2", "cat", "aaa")
	checkGetProducerInfo(t, stub, "1")
	checkGetProducerInfo(t, stub, "2")

}

func Test_Invoke2(t *testing.T) {
	scc := new(Processor)
	stub := shim.NewMockStub("test2", scc)

	checkInit(t, stub, [][]byte{[]byte("test")})

	// Invoke A->B for 123
	checkPutProducerInfo(t, stub, "1", "pig", "fusheng")
	checkPutProducerInfo(t, stub, "2", "cat", "aaa")
	checkPutProducerInfo(t, stub, "2", "egg", "aaa")
	checkPutProducerInfo(t, stub, "2", "food", "aaa")
	checkPutProducerInfo(t, stub, "2", "shop", "aaa")
	//getHistoryForKey is not supported
	// checkGetHistory(t, stub, "2")
}

func Test_Invoke3(t *testing.T) {
	scc := new(Processor)
	stub := shim.NewMockStub("test3", scc)

	checkInit(t, stub, [][]byte{[]byte("test")})

	// Invoke A->B for 123
	checkPutProducerInfo(t, stub, "1", "pig", "123")
	checkPutProducerInfo(t, stub, "2", "cat", "abcd")
	// checkGetProducerInfo(t, stub, "1")
	// getHistoryForKey is not supported
	checkGetFactoryInfo(t, stub, "2")
	// checkPutFactoryInfo(t, stub, "2", "cat", "abcd")
}
