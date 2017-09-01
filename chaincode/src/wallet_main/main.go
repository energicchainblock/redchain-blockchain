package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
	. "wallet"
)

var (
	_handleFunc map[string](func(shim.ChaincodeStubInterface, string, string, string) pb.Response)
)

func init() {
	_handleFunc = make(map[string](func(shim.ChaincodeStubInterface, string, string, string) pb.Response))
	_handleFunc["init"] = InitHandle
	_handleFunc["payment"] = TransferHandle
	_handleFunc["refund"] = TransferHandle
	_handleFunc["reward"] = RewardHandle
	_handleFunc["f-to-f"] = TransferHandle
}

type WalletChain struct{}

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
	if fun, ok := _handleFunc[cmd]; ok {
		return fun(stub, from, to, param)
	}
	return shim.Error("Invalid invoke function name.\r\n")
}

func main() {
	err := shim.Start(new(WalletChain))
	if err != nil {
		fmt.Printf("Error starting Wallet chaincode: %v\r\n", err)
	}
}
