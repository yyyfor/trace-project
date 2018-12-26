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

func checkgetInfo(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getInfo"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("getInfo", id, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("getInfo", id, "failed to get value")
		t.FailNow()
	}
	// if string(res.Payload) != value {
	// 	fmt.Println("getProduct value", id, "was not", value, "as expected")
	// 	t.FailNow()
	// }
	fmt.Println(string(res.Payload))
}

func checkPutInfo(t *testing.T, stub *shim.MockStub, id string, productName string, factory string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("putInfo"), []byte(id), []byte(productName), []byte(factory)})
	if res.Status != shim.OK {
		fmt.Println("putInfo", id, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkGetHistory(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getIdHistory"), []byte(id)})
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
	checkPutInfo(t, stub, "1", "pig", "fusheng")
	checkPutInfo(t, stub, "2", "cat", "aaa")
	checkgetInfo(t, stub, "1")
	checkgetInfo(t, stub, "2")

}

func Test_Invoke2(t *testing.T) {
	scc := new(Processor)
	stub := shim.NewMockStub("test2", scc)

	checkInit(t, stub, [][]byte{[]byte("test")})

	// Invoke A->B for 123
	checkPutInfo(t, stub, "1", "pig", "fusheng")
	checkPutInfo(t, stub, "2", "cat", "aaa")
	checkPutInfo(t, stub, "2", "egg", "aaa")
	checkPutInfo(t, stub, "2", "food", "aaa")
	checkPutInfo(t, stub, "2", "shop", "aaa")
	//getHistoryForKey is not supported
	// checkGetHistory(t, stub, "2")
}
