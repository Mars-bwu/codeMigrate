package Interaction

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	contract "github.com/ethereum/go-ethereum/contracts/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"github.com/DSiSc/crypto-suite/crypto"
	"github.com/ethereum/go-ethereum/contracts/util"
	vmutil "github.com/ethereum/go-ethereum/core/vm/util"
	wtypes "github.com/ethereum/go-ethereum/wallet/core/types"
	wutils "github.com/ethereum/go-ethereum/wallet/utils"
	wcmn "github.com/ethereum/go-ethereum/web3go/common"
)

const(
	FAILED = iota
	SUCCESS
	PENDING
)

type CrossChainPort string
// define specified type of system contract
const (
	Null = "Null"
	JustitiaChainA = "chainA"
	JustitiaChainB = "chainB"
)
const (
	InitialCrossChainPort CrossChainPort = "0"
	ChainACrossChainPort = "47768"
	ChainBCrossChainPort = "47769"
)

type Status uint64

var CrossChainAddr = vmutil.HexToAddress("0000000000000000000000000000000000011100")
var (
	//取方法名hash的前四个值
	forwardFundsMethodHash = string(util.ExtractMethodHash(util.Hash([]byte("forwardFunds(string, uint64, string)"))))
	getTxStateMethodHash = string(util.ExtractMethodHash(util.Hash([]byte("getTxState(string,string)"))))
	//为什么只有转账的没有ReceiveFunds？？
)

type CrossChainContract struct {
	records map[common.Address]Status
}

func NewCrossChainContract() *CrossChainContract {
	return new(CrossChainContract)
}

//找到要跨链的端口，并返回端口号
func CrossTargetChainPort(chainFlag string) CrossChainPort {
	var crossPort = InitialCrossChainPort
	if chainFlag == JustitiaChainA {
		crossPort = ChainACrossChainPort//如47768
	} else if chainFlag == JustitiaChainB {
		crossPort = ChainBCrossChainPort
	}
	return crossPort
}
//有啥用，什么情况下用？？？
func OppositeChainPort(chainFlag string) CrossChainPort {
	var crossPort = InitialCrossChainPort
	if chainFlag == JustitiaChainA {
		crossPort = ChainBCrossChainPort
	} else if chainFlag == JustitiaChainB {
		crossPort = ChainACrossChainPort
	}
	return crossPort
}

//如何获得合约的调用者？？？保证资金安全性
func (this *CrossChainContract) forwardFunds(toAddr common.Address, amount uint64, payload string, chainFlag string) (common.Hash, bool) {
	//调用apigateway的receiveCrossTx交易
	//参数：目的合约  数量 payload  目标链名字
	//返回值：
	from, err := GetPubliceAcccount()
	//获得公共账户地址
	//如果没有取到就报错？？？为啥会发生没有取到的情况
	if err != nil {
		return common.Hash{}, false
	}
	//这里返回的hash是什么
	hash, _,_, err := CallCrossRawTransactionReq(from, "somehangeinrpc", amount, payload, chainFlag)
	if err != nil {
		return common.Hash{}, false
	}
	return hash, false
}

func (this *CrossChainContract) getTxState(address common.Address, chainFlag string) (uint64, bool) {
	switch chainFlag {
		case "chainA":
			//123
			fmt.Println()
		default:
	}

	return SUCCESS, true
}
//跨链调用合约，在这里写入web3go调用合约的内容
func CallCrossRawTransactionReq(from common.Address, to string, amount uint64, payload string, chainFlag string) (common.Hash, common.Hash,uint64, error) {
	contract.JTMetrics.ApigatewayReceivedTx.Add(1)//**
	//amount tipOutSend1中是转出多少钱，tipInSend3是区块限制todo


	//tipoutsend0 := vmutil.HexToAddress("0x00000000000000000000000000000000outsend0")
	tipOutSend1 := "0x00000000000000000000000000000000outsend1"
	tiptocall2:="0x000000000000000000000000000000000tocall2"
	 tipInSend3 := "0x000000000000000000000000000000000insend3"
	if(to==tipOutSend1){
		targetHash,_,err:=Uout1(amount , payload , chainFlag )
		if(err!=nil){
			//错误处理
			return common.Hash{}, common.Hash{},1, err
		}
		return targetHash, common.Hash{}, 0,err
	}else if(to==tiptocall2){

		//新加的方法0526-ethcall方法
		txstate:=Uin2(amount , payload , chainFlag )
		//返回1表示交易没有成功，0表示out上交易成功


		return common.Hash{}, common.Hash{},txstate, nil

	}else if(to==tipInSend3){
		targetHash,hashValue2,err:=Uin3(amount , payload , chainFlag )
		if err != nil {
			return common.Hash{}, common.Hash{},1, nil
		}

		return common.Hash(targetHash), common.Hash(hashValue2),0, nil

	}else {
		test, _ := hex.DecodeString(payload)
		//这个可以不要---

		chainFlag_:=strings.Split(chainFlag, ":")
		web, err := wutils.NewWeb3(string(chainFlag_[0]), string(chainFlag_[1]), false)
		if err != nil {
			return common.Hash{}, common.Hash{}, 1,nil
		}
		targetHash, err := web.Eth.SendRawTransaction(test)

		return common.Hash(targetHash), common.Hash{},0, err
	}
}


//获得公共账户地址，是否可以利用托管账户保证资金安全??
//这里的公共账户指的是谁的账户
func GetPubliceAcccount() (common.Address, error){
	//get from config or genesis ?
	addr := "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"
	address := common.HexToAddress(addr)

	return common.Address(address), nil
}


//uout1-----------------------------------------


func Uout1(amount uint64, payload string, chainFlag string)(common.Hash, common.Hash, error){
	//跨链交易发出操作一

	//amount uint64 转出方转出多少钱
	//payload string 构造成的参数“签名交易数据：调用的合约地址：out中的txid”
	//chainFlag string ip:port

	//解析ip：port
	chainFlag_:=strings.Split(chainFlag, ":")
	web, err := wutils.NewWeb3(string(chainFlag_[0]), string(chainFlag_[1]), false)

	if(err!=nil){
		return common.Hash{}, common.Hash{}, err
	}
	//第一项解析payload
	payload1:=strings.Split(payload, ":")



	//out1中实际的参数
	payload1_1:=payload1[0]

	sss:=string(payload1_1)


	inContractAddress_1:=strings.ToLower(payload1[1])

	txID_:=payload1[2]
	txID1_ := hex.EncodeToString([]byte(txID_))
	txID_1:=completionStr(txID1_)

	amount_:= strconv.FormatInt(int64(amount), 16)

	amount_1:=completionuInt(amount_)
	//

	//为了支持链B对链A的转账，先不进行验证，偷个懒
	//chainFlag_1:="636861696e410000000000000000000000000000000000000000000000000000"
	//本链标识chainA的16进制string补全值

	inContractAddress_2 ,outAmount_2 ,txID_2,isOk:= get4Param(sss)
	if(isOk==true){

		return common.Hash{}, common.Hash{}, err
	}


	if(inContractAddress_2!=inContractAddress_1 ||outAmount_2!=amount_1 ||txID_2!=txID_1){

		return common.Hash{}, common.Hash{}, err

	}

	test, _ := hex.DecodeString(payload1_1)

	targetHash, err := web.Eth.SendRawTransaction(test)
	if err != nil {
		//错误处理
	}

	return common.Hash(targetHash), common.Hash{}, err

}

func get4Param(payload string )(string ,string,string,bool ){
	//合约地址，交易ID，out减去的资金，chainA类似

	test, _ := hex.DecodeString(payload)

	tx := new(contract.Transaction)
	if err := rlp.DecodeBytes(test, tx); err != nil {

		ethTx := new(contract.ETransaction)
		err = ethTx.DecodeBytes(test)
		if err != nil {

			//错误处理
			return "","","",true
		}
		ethTx.SetTxData(&tx.Data)
	}

	//实际调用的合约地址
	txToAddr := types.AddressToHex(*tx.Data.Recipient)

	//得到的实际input
	input := tx.Data.Payload

	rawInput:= util.BytesToHex(input)


	if(len(rawInput)!=650){

		return "","","",true
	}

	outAmount_:=rawInput[74:138]//outAmount

	txID_:=rawInput[458:522]//txID

	//chainFlag_:=rawInput[586:650]//chainFlag


	return txToAddr,string(outAmount_),string(txID_),false
}



//-----------------------------------------uout1


//in2---------------------------------------

func Uin2(limit uint64, payload string, chainFlag string)(uint64){
	//tip, limit, payload,chainFlag
	//payload=txID:outChainContractAddress



	//第一项解析payload中的txID，outChainContractAddress
	payload1:=strings.Split(payload, ":")
	txID:=string(payload1[0])
	outChainContractAddress:=string(payload1[1])
	from:="0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"

	txstate:=checkChainTxStateSecond(chainFlag,txID,outChainContractAddress,from)

	return uint64(txstate)

}

func checkChainTxStateSecond(ip string,txid string,contractAddress string,accountAddress string )(int){

	url:="http://"+ip

	//构建用于发送的json
	request :=createDataJson(txid ,contractAddress,accountAddress)
	byteData,_ := json.Marshal(&request)
	reader := bytes.NewReader(byteData)
	for i := 0; i <= 3; i++ {

		resp, err := http.Post(url, "application/x-www-form-urlencoded", reader)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 1
		}

		var recipient RequestCall
		json.Unmarshal(body, &recipient)

		//得到返回的状态结果
		aa := recipient.Result
		bb := "0x0000000000000000000000000000000000000000000000000000000000000002"

		if aa == bb {
			return 0
		}
		time.Sleep(time.Second * 3)

	}
	return 1
	//0正确1有误

}

func createDataJson(txid string,contractAddress string,accountAddress string)CallJson{

	txID1:= hex.EncodeToString([]byte(txid ))
	txID2:=completionStr1(txID1)
	param1:="0x3f006398"
	param2:="00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012"
	data:=param1+param2+txID2


	var aa MyData=Data{accountAddress ,contractAddress,data,"0x0","0x2dc6c0"}

	var bb CallJson
	bb.Params[0]=aa
	bb.Params[1]="latest"
	bb.Jsonrpc="2.0"
	bb.ID=1
	bb.Method="eth_call"

	fmt.Print(bb)
	return bb
}

//发送的第二层
type CallJson struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  [2]interface{} `json:"params"`
}

//发送的第一层与第二层之间的连接
type MyData interface{

}

//发送的第一层
type Data struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Data  string `json:"data"`
	Value string `json:"value"`
	Gas   string `json:"gas"`
}

//查到的结果
type RequestCall struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}



//后面补零成64位
func completionStr1(value string)string {

	a:=int(len(value))
	rawData :="0000000000000000000000000000000000000000000000000000000000000000"
	rawData_:=string(rawData[a:])
	rawData_1:=value+rawData_
	return rawData_1

}




//-----------------------------------------in2




//in3---------------------------------------

func Uin3(limit uint64, payload string, chainFlag string)(common.Hash, common.Hash, error){
	//tip, limit, payload,chainFlag
	//payload=txID:outChainContractAddress

	chainFlag_:=strings.Split(chainFlag, ":")
	web, err := wutils.NewWeb3(string(chainFlag_[0]), string(chainFlag_[1]), false)
	if err != nil {
		return common.Hash{}, common.Hash{}, nil
	}
	//第一项解析payload中的txID，outChainContractAddress
	payload1:=strings.Split(payload, ":")
	txID:=string(payload1[0])

	createInput,haveError:=createInput(txID)
	if(haveError==true){
		return common.Hash{}, common.Hash{}, nil
	}

	input_ := wcmn.HexToBytes(createInput)
	outChainContractAddress:=string(payload1[1])
	to := common.HexToAddress(outChainContractAddress)


	from:=common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B")
	bigNonce , err := web.Eth.GetTransactionCount(wcmn.Address(from), "latest")
	if err != nil {
		return common.Hash{}, common.Hash{}, err
	}

	tx_ := new(contract.Transaction)
	addr, err := GetPubliceAcccount2()
	//这个公用账户必须在所有链上都有
	if err != nil {
		return common.Hash{}, common.Hash{}, nil
	}
	tx_.Data.From = &addr

	//private := "29ad43a4ebb4a65436d9fb116d471d96516b3d5cc153e045b384664bed5371b9"
	private := "45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8"
	nonce := bigNonce.Uint64()
	tx_.Data.AccountNonce = nonce
	tx_.Data.Price = big.NewInt(0)
	tx_.Data.GasLimit = 6721975
	tx_.Data.Recipient = &to
	tx_.Data.Amount = big.NewInt(int64(0))
	//payload填充问题，填充合约调用的参数
	tx_.Data.Payload = input_

	//sign tx
	priKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return common.Hash{}, common.Hash{}, nil
	}

	chainID := big.NewInt(int64(1))
	tx_, err = wtypes.SignTx(tx_, wtypes.NewEIP155Signer(chainID), priKey)
	if err != nil {
		return common.Hash{}, common.Hash{}, nil
	}

	txBytes, err:= rlp.EncodeToBytes(tx_)
	if err != nil {
		return common.Hash{}, common.Hash{}, nil
	}

	//这个可以不要---写这个只是为了得到签完名的交易数据
	//encodedStr := hex.EncodeToString(txBytes)

	//test, _ := hex.DecodeString(encodedStr)
	//这个可以不要---

	targetHash, err := web.Eth.SendRawTransaction(txBytes)
	if err != nil {
		return common.Hash{}, common.Hash{}, nil
	}
	return common.Hash(targetHash), common.Hash(targetHash), nil

}

func GetPubliceAcccount2() (common.Address, error){
	//get from config or genesis ?
	addr := "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"
	address := common.HexToAddress(addr)

	return common.Address(address), nil
}


func createInput(txID string )(string,bool){

	functionName:="0xe4b808be00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012"
	if(len(txID)!=18){

		return "",true

	}

	txID1_ := hex.EncodeToString([]byte(txID))
	txID_1:=completionStr(txID1_)

	input:=functionName+txID_1

	return input,false
}



//-----------------------------------------in3


//前面补零成64位
func completionuInt(value string)string {

	a:=int(len(value))
	rawData :="0000000000000000000000000000000000000000000000000000000000000000"

	rawData_:=string(rawData[a:])
	rawData_1:=rawData_+value
	return rawData_1
}
//后面补零成64位
func completionStr(value string)string {

	a:=int(len(value))
	rawData :="0000000000000000000000000000000000000000000000000000000000000000"
	rawData_:=string(rawData[a:])
	rawData_1:=value+rawData_
	return rawData_1

}


