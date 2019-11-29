package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type ChenTaoCC01 struct {
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (cc *ChenTaoCC01) Init(stub shim.ChaincodeStubInterface) pb.Response {
	chanID := stub.GetChannelID()
	fmt.Println("chentaocc01 init with channel ID: ", chanID)

	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (cc *ChenTaoCC01) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println("chentaocc01 invoke with channel ID: ", stub.GetChannelID())

	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "create":
		return cc.create(stub, args)

	case "delete":
		return cc.delete(stub, args)

	case "deposit":
		return cc.deposit(stub, args)

	case "withdraw":
		return cc.withdraw(stub, args)

	case "pay":
		return cc.pay(stub, args)

	}

	return shim.Error("Invalid invoke function name. Expecting \"create\" \"delete\" \"deposit\" \"withdraw\" \"pay\" ")
}

func (cc *ChenTaoCC01) debug(flag, channelID, uid string, val int) {
	fmt.Printf("[%s chentaocc01 channelID=%s] UID = %s, VALUE = %d \n", flag, channelID, uid, val)
}

func (cc *ChenTaoCC01) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	uid := args[0]
	val, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("expecting integer value for account asset")
	}

	cc.debug("CREATE", stub.GetChannelID(), uid, val)

	err = stub.PutState(uid, []byte(strconv.Itoa(val)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cc *ChenTaoCC01) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("incorrect count of arguments")
	}

	uid := args[0]

	cc.debug("DELETE", stub.GetChannelID(), uid, 0)

	if err := stub.DelState(uid); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cc *ChenTaoCC01) deposit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	uid := args[0]
	val, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("expecting integer value for account asset")
	}

	cc.debug("DEPOSIT", stub.GetChannelID(), uid, val)

	// get user's current balance
	valBytes, err := stub.GetState(uid)
	if err != nil {
		return shim.Error(err.Error())
	}

	if valBytes == nil {
		return shim.Error("no account")
	}

	currentVal, err := strconv.Atoi(string(valBytes))
	if err != nil {
		return shim.Error(err.Error())
	}

	currentVal = currentVal + val
	err = stub.PutState(uid, []byte(strconv.Itoa(currentVal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cc *ChenTaoCC01) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	uid := args[0]
	cc.debug("QUERY", stub.GetChannelID(), uid, 0)

	// get user's current balance
	valBytes, err := stub.GetState(uid)
	if err != nil {
		return shim.Error(err.Error())
	}

	if valBytes == nil {
		return shim.Error("no account")
	}

	fmt.Printf("Query Response:%s\n", string(valBytes))

	return shim.Success(valBytes)
}

func (cc *ChenTaoCC01) withdraw(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	uid := args[0]
	val, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("expecting integer value for account asset")
	}

	cc.debug("WITHDRAW", stub.GetChannelID(), uid, val)

	// get user's current balance
	valBytes, err := stub.GetState(uid)
	if err != nil {
		return shim.Error(err.Error())
	}

	if valBytes == nil {
		return shim.Error("no account")
	}

	currentVal, err := strconv.Atoi(string(valBytes))
	if err != nil {
		return shim.Error(err.Error())
	}

	if currentVal < val {
		return shim.Error("not enough balance")
	}

	currentVal = currentVal - val
	err = stub.PutState(uid, []byte(strconv.Itoa(currentVal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

func (cc *ChenTaoCC01) pay(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	uidFrom := args[0]
	uidTo := args[1]
	val, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("expecting integer value for account asset")
	}

	cc.debug("PAY", stub.GetChannelID(), "["+uidFrom+"->"+uidTo+"]", val)

	// get user's current balance
	valFromBytes, err := stub.GetState(uidFrom)
	if err != nil {
		return shim.Error(err.Error())
	}

	if valFromBytes == nil {
		return shim.Error("no from account")
	}

	currentValFrom, err := strconv.Atoi(string(valFromBytes))
	if err != nil {
		return shim.Error(err.Error())
	}

	valToBytes, err := stub.GetState(uidTo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if valToBytes == nil {
		return shim.Error("no to account")
	}

	currentValTo, err := strconv.Atoi(string(valToBytes))
	if err != nil {
		return shim.Error(err.Error())
	}

	if currentValFrom < val {
		return shim.Error("not enough from balance")
	}

	currentValFrom = currentValFrom - val
	err = stub.PutState(uidFrom, []byte(strconv.Itoa(currentValFrom)))
	if err != nil {
		return shim.Error(err.Error())
	}

	currentValTo = currentValTo + val
	err = stub.PutState(uidTo, []byte(strconv.Itoa(currentValTo)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {

	err := shim.Start(new(ChenTaoCC01))
	if err != nil {
		fmt.Printf("%s", err)
	}

}
