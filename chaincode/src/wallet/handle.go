package wallet

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// func KeysHandle(stub shim.ChaincodeStubInterface, from string, to string, param string) pb.Response {
// 	req := &KeysReq{}
// 	err := json.Unmarshal([]byte(param), req)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	keys, err := stub.GetStateByRange(req.Start, req.End)
// 	for {
// 		if keys.HasNext() {
// 			kv, _ := keys.Next()
// 			fmt.Printf("kv=%v\r\n", kv)
// 		} else {
// 			break
// 		}
// 	}
// 	keys.Close()
// 	return shim.Success(nil)
// }

func InitHandle(stub shim.ChaincodeStubInterface, from, to, param string) pb.Response {
	req := &Wallet{}
	err := json.Unmarshal([]byte(param), req)
	if err != nil {
		return shim.Error(err.Error())
	}
	if req.Available < 0 {
		return shim.Error("available cannot be less than zero")
	}
	// if req.Ico < 0 {
	// 	return shim.Error("ico cannot be less than zero")
	// }
	toBytes, err := stub.GetState(to)
	if err != nil {
		return shim.Error(err.Error())
	}
	if toBytes != nil {
		return shim.Error("addr already exists")
	}
	err = stub.PutState(to, []byte(param))
	if err != nil {
		return shim.Error(err.Error())
	}
	// retList := make([]*CommonReply, 0)
	// retList = append(retList, &CommonReply{
	// 	Address:   to,
	// 	Available: req.Available,
	// })
	reply := &CommonReply{
		Address: to,
		Value:   req.Available,
		From:    "",
		Fvalue:  0,
		Amount:  req.Available,
	}
	ret, _ := json.Marshal(reply)
	return shim.Success(ret)
}

func TransferHandle(stub shim.ChaincodeStubInterface, from, to, param string) pb.Response {
	req := &CommonReq{}
	err := json.Unmarshal([]byte(param), req)
	if err != nil {
		return shim.Error(err.Error())
	}
	if req.Number < 0 {
		return shim.Error("number cannot be less than zero")
	}
	fromBytes, err := stub.GetState(from)
	if err != nil {
		return shim.Error(err.Error())
	}
	if fromBytes == nil {
		return shim.Error("from addr not exist")
	}
	toByte, err := stub.GetState(to)
	if err != nil {
		return shim.Error(err.Error())
	}
	//if toByte == nil {
	//	return shim.Error("to addr not exist")
	//}
	fromWallet := &Wallet{}
	err = json.Unmarshal(fromBytes, fromWallet)
	if err != nil {
		return shim.Error(err.Error())
	}
	toWallet := &Wallet{}
	if toByte != nil {
		err = json.Unmarshal(toByte, toWallet)
		if err != nil {
			return shim.Error(err.Error())
		}
	} else {
		toWallet.Available = 0
	}
	if fromWallet.Available < req.Number {
		return shim.Error("from addr have not enough coin")
	}
	fromWallet.Available -= req.Number
	toWallet.Available += req.Number
	fromState, _ := json.Marshal(fromWallet)
	err = stub.PutState(from, fromState)
	if err != nil {
		return shim.Error(err.Error())
	}
	toState, _ := json.Marshal(toWallet)
	err = stub.PutState(to, toState)
	if err != nil {
		return shim.Error(err.Error())
	}
	// retList := make([]*CommonReply, 0)
	// retList = append(retList, &CommonReply{
	// 	Address:   from,
	// 	Available: fromWallet.Available,
	// })
	// retList = append(retList, &CommonReply{
	// 	Address:   to,
	// 	Available: toWallet.Available,
	// })
	reply := &CommonReply{
		Address: to,
		Value:   toWallet.Available,
		From:    from,
		Fvalue:  fromWallet.Available,
		Amount:  req.Number,
	}
	ret, _ := json.Marshal(reply)
	return shim.Success(ret)
}

func RewardHandle(stub shim.ChaincodeStubInterface, from, to, param string) pb.Response {
	req := &CommonReq{}
	err := json.Unmarshal([]byte(param), req)
	if err != nil {
		return shim.Error(err.Error())
	}
	if req.Number < 0 {
		return shim.Error("number cannot be less than zero")
	}
	toBytes, err := stub.GetState(to)
	if err != nil {
		return shim.Error(err.Error())
	}
	if toBytes == nil {
		return shim.Error("to addr not exist")
	}
	toWallet := &Wallet{}
	err = json.Unmarshal(toBytes, toWallet)
	if err != nil {
		return shim.Error(err.Error())
	}
	toWallet.Available += req.Number
	toState, _ := json.Marshal(toWallet)
	err = stub.PutState(to, toState)
	if err != nil {
		return shim.Error(err.Error())
	}
	// retList := make([]*CommonReply, 0)
	// retList = append(retList, &CommonReply{
	// 	Address:   to,
	// 	Available: toWallet.Available,
	// })
	reply := &CommonReply{
		Address: to,
		Value:   toWallet.Available,
		From:    from,
		Fvalue:  0,
		Amount:  req.Number,
	}
	ret, _ := json.Marshal(reply)
	return shim.Success(ret)
}

func QueryHandle(stub shim.ChaincodeStubInterface, key string) pb.Response {
	stateBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("addr not exist")
	}
	return shim.Success(stateBytes)
}

func TotalHandle(stub shim.ChaincodeStubInterface, q string) pb.Response {
	keys, err := stub.GetQueryResult(q)
	if err != nil {
		return shim.Error(err.Error())
	}
	for {
		if keys.HasNext() {
			kv, _ := keys.Next()
			fmt.Printf("kv=%v\r\n", kv)
		} else {
			break
		}
	}
	keys.Close()
	return shim.Success(nil)
}
