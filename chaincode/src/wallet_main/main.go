package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
	. "wallet"
)

type WalletInterface interface {
	Wallet(stub shim.ChaincodeStubInterface, from string, to string, param string) pb.Response
	Recharge(stub shim.ChaincodeStubInterface, from string, to string, param string) pb.Response
	Transfer(stub shim.ChaincodeStubInterface, from string, to string, param string) pb.Response
}

type WalletChain struct{}

var (
	_handle WalletInterface
)

func init() {

}

func (w *WalletChain) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (w *WalletChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	fmt.Printf("args=%v\r\n", args)
	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments.")
	}
	subFunc := args[0]
	cmd := args[1]
	to := args[2]
	from := args[3]
	param := args[4]
	fmt.Printf("subfunc=%v, cmd=%v, to=%v, from=%v, param=%v\r\n", subFunc, cmd, to, from, param)
	cmd = strings.ToLower(cmd)
	if cmd == "wallet" {
		return _handle.Wallet(stub, from, to, param)
	} else if cmd == "recharge" {
		return _handle.Recharge(stub, from, to, param)
	} else if cmd == "transfer" {
		return _handle.Transfer(stub, from, to, param)
	}
	return shim.Error("Invalid invoke function name.\r\n")
}

func main() {
	_handle = &WalletHandle{}
	err := shim.Start(new(WalletChain))
	if err != nil {
		fmt.Printf("Error starting Wallet chaincode: %v\r\n", err)
	}
}
