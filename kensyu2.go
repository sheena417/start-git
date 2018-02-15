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
  var initOwner = Owner{ObjectType:"Owner",Id:"owner",Quantity:"0"}
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
  var bank = Bank{ObjectType:"Bank",BankCode:args[0]}//[]Account
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

  balance64, err := strconv.ParseFloat(args[3], 64)
  if err != nil {
    return shim.Error("Incorrect arguments. Expecting a UserID and balance.")
  }
  if balance64 < 0  {
    return shim.Error("Incorrect arguments. Expecting a UserID and balance.")
  }

  var accounts = Account{ObjectType:"Account",AccountNumber:args[0],Id:args[1],BankCode:args[2],Balance:balance64}
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
//********************************************************//
  txid, err := APIstub.GetState(args[0])
  if err != nil {
      return shim.Error(fmt.Sprintf("Already exists txid: %s with error: %s", args[0], err))
  }

  if txid != nil {
    return shim.Error(fmt.Sprintf("Txid already exist: %s with error: %s", args[0], err))
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

  quantity64, err := strconv.ParseFloat(args[3], 64)
  if err != nil {
      return shim.Error("Incorrect arguments. Expecting a UserID and balance.")
  }

  if quantity64 < 0 {
    return shim.Error("Incorrect arguments. Expecting a balance.")
  }

  fee64, err := strconv.ParseFloat(args[4], 64)
  if err != nil {
      return shim.Error("Incorrect arguments. Expecting a fee.")
  }
  if fee64 < 0 {
    return shim.Error("Incorrect arguments. Expecting a fee.")
  }
//********************************************************//

  //　Fromの残高を確認する
  fromAccountItf := Account{}
  json.Unmarshal(fromAccount,&fromAccountItf)
  var feeBalance = quantity64 + fee64
  if fromAccountItf.Balance < feeBalance {
      return shim.Error("Not enough fromUsers balance.")
  }

  // ToをUnmarshal
  toAccountItf := Account{}
  json.Unmarshal(toAccount,&toAccountItf)
  // 残高を加算
  toAccountItf.Balance = toAccountItf.Balance + quantity64
  //残高を減算
  fromAccountItf.Balance = fromAccountItf.Balance - quantity64

  //FromAccountの残高を反映
  fromAccountBytes,_ := json.Marshal(fromAccountItf)
  APIstub.PutState(args[1],fromAccountBytes)

  //ToAccountの残高を反映
  toAccountBytes,_ := json.Marshal(toAccountItf)
  APIstub.PutState(args[2],toAccountBytes)

  //add fee to owner account
  owner, err := APIstub.GetState("owner")
  ownerItf := Owner{}
  json.Unmarshal(owner,&ownerItf)
  ownerItf.Quantity = ownerItf.Quantity + args[4]
  ownerBytes,_ := json.Marshal(ownerItf)
  APIstub.PutState("owner",ownerBytes)

  // 最終的にTXIDをキーにして、PUTSTATEする
  var transfer = Transfer {
            ObjectType:"Transfer",
            TxId:args[0],
            FromAccount:args[1],
            ToAccount:args[2],
            Quantity:quantity64,
            Fee:fee64 }
  transferBytes,_ := json.Marshal(transfer)
  APIstub.PutState(args[0],transferBytes)

  fmt.Println("******* end transfer *******")
  return shim.Success(nil)

}

//
func (c *Chain) query(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
  fmt.Println("******* start query *******")
  if len(args) != 1 {
      return shim.Error("Incorrect arguments. Expecting a Particular Key")
  }

  query, err := APIstub.GetState(args[0])
  if err != nil {
      return shim.Error(fmt.Sprintf("Not exists KeyValue: %s with error: %s", args[0], err))
  }
  if query == nil{
      return shim.Error(fmt.Sprintf("No Value this KeyValue: %s", args[0]))
  }

  fmt.Println(query)


  fmt.Println("******* end query *******")
  return shim.Success(nil)
}




func main() {
  err := shim.Start(new(Chain));
  if err != nil {
      fmt.Printf("Error starting Chain chaincode: %s", err)
  }
}
