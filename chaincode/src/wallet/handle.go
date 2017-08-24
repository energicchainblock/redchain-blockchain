package wallet

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type WalletHandle struct{}

func (wh *WalletHandle) Wallet(stub shim.ChaincodeStubInterface, from string, to string, param string) pb.Response {
	addr := from
	walletBytes, err := stub.GetState(addr)
	if err != nil {
		return shim.Error(err.Error())
	}
	wallet := &Wallet{}
	if walletBytes == nil {
		wallet.Balance = float64(0.0)
		wallet.In = float64(0.0)
		wallet.Out = float64(0.0)
	} else {
		err = json.Unmarshal(walletBytes, wallet)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	ret, _ := json.Marshal(wallet)
	return shim.Success(ret)
}

func (wh *WalletHandle) Recharge(stub shim.ChaincodeStubInterface, from string, to string, param string) pb.Response {
	addr := from
	walletBytes, err := stub.GetState(addr)
	if err != nil {
		return shim.Error(err.Error())
	}
	wallet := &Wallet{}
	if walletBytes == nil {
		wallet.Balance = float64(0.0)
		wallet.In = float64(0.0)
		wallet.Out = float64(0.0)
	} else {
		err = json.Unmarshal(walletBytes, wallet)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	req := &RechargeReq{}
	err = json.Unmarshal([]byte(param), req)
	if err != nil {
		return shim.Error(err.Error())
	}
	wallet.Balance += req.Value
	wallet.In += req.Value
	ret, _ := json.Marshal(wallet)
	err = stub.PutState(addr, ret)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(ret)
}

func (wh *WalletHandle) Transfer(stub shim.ChaincodeStubInterface, from string, to string, param string) pb.Response {
	req := &TransferReq{}
	err := json.Unmarshal([]byte(param), req)
	if err != nil {
		return shim.Error(err.Error())
	}
	fromBytes, err := stub.GetState(from)
	if err != nil {
		return shim.Error(err.Error())
	}
	fromWallet := &Wallet{}
	if fromBytes == nil {
		fromWallet.Balance = float64(0.0)
		fromWallet.In = float64(0.0)
		fromWallet.Out = float64(0.0)
	} else {
		err = json.Unmarshal(fromBytes, fromWallet)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	toBytes, err := stub.GetState(to)
	if err != nil {
		return shim.Error(err.Error())
	}
	toWallet := &Wallet{}
	if toBytes == nil {
		toWallet.Balance = float64(0.0)
		toWallet.In = float64(0.0)
		toWallet.Out = float64(0.0)
	} else {
		err = json.Unmarshal(toBytes, toWallet)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	if fromWallet.Balance < req.Value {
		return shim.Error("not enough balance")
	}

	fromWallet.Balance -= req.Value
	fromWallet.Out += req.Value
	toWallet.Balance += req.Value
	toWallet.In += req.Value

	retFrom, _ := json.Marshal(fromWallet)
	retTo, _ := json.Marshal(toWallet)

	err = stub.PutState(from, retFrom)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(to, retTo)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
