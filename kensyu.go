package main

import (
    "encoding/json"
    "fmt"
    "strconv"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)



type Chain struct {
}

type Owner struct {
  ObjectType string `json:docType`
  Id string `json:id`
  Quantity string `json:quantity`
}

type Bank struct {
  ObjectType string `json:docType`
  BankCode string `json:bankCode`
  Accounts []string `json:accounts`
}

type Account struct {
  ObjectType string `json:docType`
  AccountNumber string `json:accountNumber`
  Id string `json:id`
  BankCode string `json:bankCode`
  Balance float64 `json:balance`
}

type Transfer struct {
  ObjectType string `json:docType`
  TxId string `json:txId`
  FromAccount string `json:fromAccount`
  ToAccount string `json:toAccount`
  Quantity float64 `json:quantity`
  Fee float64 `json:fee`
}


//INIT method
func (c *Chain) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
  fmt.Println("******* start Init *******")

  //Create Owner accounts
  initOwner := Owner{ObjectType:"Owner",Id:"owner",Quantity:"0"}
  ownerBytes,_ := json.Marshal(initOwner)
  APIstub.PutState(initOwner.Id , ownerBytes)

  fmt.Println("******* end Init *******")
  return shim.Success(nil)
}

//Invoke method
func (c *Chain) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
  //function:
  //args:
  function, args := APIstub.GetFunctionAndParameters()

  //createBank
  if function == "createBank" {
      fmt.Println("******* go createBank *******")
      return c.createBank(APIstub, args)
  //createAccount
  } else if function == "createAccount" {
      fmt.Println("******* go createAccount *******")
      return c.createAccount(APIstub, args)
  //transfer
  } else if function == "transfer" {
      fmt.Println("******* go transfer *******")
      return c.transfer(APIstub, args)
  //query
  } else if function == "query" {
      fmt.Println("******* go query *******")
      return c.query(APIstub, args)
  }
  return shim.Error("Invalid function name.")
}




//
func (c *Chain) createBank(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
  fmt.Println("******* start createBank *******")

  if len(args) != 1 {
      return shim.Error("Incorrect arguments. Expecting a bankCode")
  }
  var bank = Bank{ObjectType:"Bank",BankCode:args[0]}
  bankBytes,_ := json.Marshal(bank)
  APIstub.PutState(args[0],bankBytes)

  fmt.Println("******* end createBank *******")

  return shim.Success(nil)
  }

//
func (c *Chain) createAccount(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
  fmt.Println("******* start createAccount *******")
  if len(args) != 4 {
  return shim.Error("Incorrect arguments. Expecting a AccountNumber, Id ,bc and balance")
  }

  account, err := APIstub.GetState(args[0])
  if err != nil {
      return shim.Error(fmt.Sprintf("Already exist this AccountNumber: %s with error: %s", args[0], err))
  }
  if account == nil {
    return shim.Error(fmt.Sprintf("Already exist this AccountNumber: %s with error: %s", args[0], err))
  }


  id, err := APIstub.GetState(args[1])
  if err != nil {
    return shim.Error(fmt.Sprintf("Already exist this: %s with error: %s", args[1], err))
  }
  if id == nil {
    return shim.Error(fmt.Sprintf("Already exist this: %s with error: %s", args[1], err))
  }

  bankCode, err := APIstub.GetState(args[2])
  if err == nil {
    return shim.Error(fmt.Sprintf("Not exist this Bank: %s with error: %s", args[2], err))
  }
  if bankCode == nil {
    return shim.Error(fmt.Sprintf("Not exist this Bank: %s with error: %s", args[2], err))
  }

  balance, err := strconv.ParseFloat(args[3], 32)
  if err != nil {
    return shim.Error("Incorrect arguments. Expecting a UserID and balance.")
  }
  if balance < 0  {
    return shim.Error("Incorrect arguments. Expecting a UserID and balance.")
  }

  var accounts = Account{ObjectType:"Account",AccountNumber:args[0],Id:args[1],BankCode:args[2],Balance:balance}
  accountBytes,_ := json.Marshal(accounts)
  APIstub.PutState(args[0],accountBytes)

  //Bank内の該当の銀行コードの[]Account配列の中にこの構造体（Account）を入れたい



  fmt.Println("******* end createAccount *******")
  return shim.Success(nil)
}

//
func (c *Chain) transfer(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
  fmt.Println("******* start transfer *******")
  if len(args) != 5 {
      return shim.Error("Incorrect arguments. Expecting a TxId, FromAccountNumber, ToAccountNumber, Quantity and Fee")
  }
  //key, value でgetstate
  txid64, err := APIstub.GetState(args[0])
  if err != nil {
      return shim.Error(fmt.Sprintf("Failed to get txid: %s with error: %s", args[0], err))
  }

  if txid64 == nil {
    return shim.Error(fmt.Sprintf("Failed to get txid: %s with error: %s", args[0], err))
  }

  fromAccount, err := APIstub.GetState(args[1])
  if err != nil {
      return shim.Error(fmt.Sprintf("Failed to get balance: %s with error: %s", args[1], err))
  }
  if fromAccount == nil{
      return shim.Error(fmt.Sprintf("FromUser not found: %s", args[1]))
  }

  toAccount, err := APIstub.GetState(args[2])
  if err != nil {
      return shim.Error(fmt.Sprintf("Failed to get balance: %s with error: %s", args[2], err))
  }
  if toAccount == nil{
      return shim.Error(fmt.Sprintf("ToUser not found: %s", args[2]))
  }

  quantity32, err := strconv.ParseFloat(args[3], 32)
  if err != nil {
      return shim.Error("Incorrect arguments. Expecting a UserID and balance.")
  }

  if quantity32 < 0 {
    return shim.Error("Incorrect arguments. Expecting a balance.")
  }

  fee32, err := strconv.ParseFloat(args[4], 32)
  if err != nil {
      return shim.Error("Incorrect arguments. Expecting a UserID and balance fee.")
  }
  if fee32 < 0 {
    return shim.Error("Incorrect arguments. Expecting a fee.")
  }


  //add fee to owner account

  fmt.Println("******* end transfer *******")
  return shim.Success(nil)

}

//
func (c *Chain) query(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
  fmt.Println("******* start query *******")
  if len(args) != 1 {
      return shim.Error("Incorrect arguments. Expecting a Particular Key")
  }





  fmt.Println("******* end query *******")
  return shim.Success(nil)
}




func main() {
  err := shim.Start(new(Chain));
  if err != nil {
      fmt.Printf("Error starting Chain chaincode: %s", err)
  }
}
