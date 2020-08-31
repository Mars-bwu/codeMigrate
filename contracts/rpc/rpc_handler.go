package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	craft "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/Interaction"
	transaction "github.com/ethereum/go-ethereum/contracts/types"
	"github.com/ethereum/go-ethereum/contracts/util"
	"github.com/ethereum/go-ethereum/core/types"
	cutil "github.com/ethereum/go-ethereum/core/vm/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	wtypes "github.com/ethereum/go-ethereum/wallet/core/types"
	wutils "github.com/ethereum/go-ethereum/wallet/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"reflect"
	"strings"

	//ctypes "github.com/DSiSc/craft/types"
)


type RecipientData struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		BlockHash         string      `json:"blockHash"`
		BlockNumber       string      `json:"blockNumber"`
		TransactionHash   string      `json:"transactionHash"`
		TransactionIndex  string      `json:"transactionIndex"`
		From              string      `json:"from"`
		To                string      `json:"to"`
		Root              string      `json:"root"`
		Status            string      `json:"status"`
		GasUsed           string      `json:"gasUsed"`
		CumulativeGasUsed string      `json:"cumulativeGasUsed"`
		LogsBloom         interface{} `json:"logsBloom"`
		Logs              []struct {
			Address          string   `json:"address"`
			Topics           []string `json:"topics"`
			Data             string   `json:"data"`
			BlockNumber      string   `json:"blockNumber"`
			TransactionHash  string   `json:"transactionHash"`
			TransactionIndex string   `json:"transactionIndex"`
			BlockHash        string   `json:"blockHash"`
			LogIndex         string   `json:"logIndex"`
			Removed          bool     `json:"removed"`
		} `json:"logs"`
		ContractAddress interface{} `json:"contractAddress"`
	} `json:"result"`
}








var RpcContractAddr = cutil.HexToAddress("0000000000000000000000000000000000011101")




// rpc routes
var routes = map[string]*RPCFunc{//表明一个变量是指针类型
	string(util.ExtractMethodHash(util.Hash([]byte("ForwardFunds(string,uint64,string,string)")))): NewRPCFunc(ForwardFunds),
	string(util.ExtractMethodHash(util.Hash([]byte("GetTxState(string,uint64,string,string)")))): NewRPCFunc(GetTxState),
	string(util.ExtractMethodHash(util.Hash([]byte("ReceiveFunds(address,uint64,string,uint64)")))): NewRPCFunc(ReceiveFunds),
}

// 0 means failed, 1 means success
//需要在这里面写入web3go或者用http请求调用Interaction.CallCrossRawTransaction实现
//超级节点要部署两条链，这会不会压力太大，但是联盟链无所谓
func ForwardFunds(toAddr string, amount uint64, payload string, chainFlag string) (error, string, string, uint64) {
	from, _ := Interaction.GetPubliceAcccount()//**是想获得谁的公共账户
	//**对方链的某个地址么,用引的一个包对他进行了处理
	localHash, targetHash,txstate, err := Interaction.CallCrossRawTransactionReq(from, toAddr, amount, payload, chainFlag)
	//**这里面是具体的Web3go的方法么？这个方法的作用？只适用于转账么？
	//这里面写的看不明白，这是web3go调用合约的写法么，看点啥能更好的理解这块的代码
	if err != nil {
		return err, "", "", txstate
	}

	localBytes := cutil.HashToBytes(localHash)//**
	targetBytes := cutil.HashToBytes(targetHash)//**

	return err, util.BytesToHex(localBytes), util.BytesToHex(targetBytes),txstate
}

// GetCross Tx state
func GetTxState(txHash string, amount uint64, tmp string, chainFlag string) (error, uint64){

	//tip用于辨识调用哪个方法
	tipOut2:="0x0000000000000000000000000000000outquery2"
	//state 1不明白 2成功  3失败

	tipIn4:="0x00000000000000000000000000000000inquery4"
	limit :=amount
	//state 1不明白  2成功  3失败


	tip:=tmp;
	if(tip==tipOut2){

		err,state:=out2(txHash,limit,chainFlag)

		return err,state
	}else if(tip==tipIn4){

		err,state:=In4(txHash,limit,chainFlag)

		return err,state
	}
	err,state:=generate(txHash ,chainFlag )
	return err,state

}

// Receipt funds
//没有进行跨链交互么
func ReceiveFunds(to string, amount uint64, payload string, srcChainId uint64) (error, uint64){
	//receive tx bytes, decode input
	input := util.HexToBytes(payload)
	tx := new(transaction.Transaction)
	//没看见有给tx赋值的地方呀
	// new用来分配内存
	//**craft.Transaction？？一个结构体
	if err := rlp.DecodeBytes(input, tx); err != nil {
		//RLP的唯一目标就是解决结构体的编码问题
		//https://www.cnblogs.com/baizx/p/6928622.html
		ethTx := new(transaction.ETransaction)
		//tx := new(craft.Transaction)与ethTx := new(craft.ETransaction)----craft.Transaction与craft.ETransaction有啥区别？？？
		//**craft.ETransaction
		//ETransaction与Transaction有啥区别
		err = ethTx.DecodeBytes(input)
		//里面也是rlp.DecodeBytes
		//err不为空就报错
		if err != nil {
			log.Info("sendRawTransaction tx decode as ethereum error, err = %v", err)
			return err, 0
		}
		ethTx.SetTxData(&tx.Data)
		//这是在干什么
	}
	//tx该结构体里面到底是什么，里面的各个值是什么意思
	from_, err := wtypes.Sender(wtypes.NewEIP155Signer(big.NewInt(int64(srcChainId))), tx)
	//wtypes.Sender得到发送者，但是里面干了什么操作不知道
	if err != nil {
		log.Error("get from address failed, err = %v", err)
		return err, 0
	}

	contractAddr := "0x47c5e40890bce4a473a49d7501808b9633f29782"
	//这个合约地址是谁的？？？？？
	//**??
	targetToInput := util.BytesToHex(tx.Data.Payload)
	txToAddr := types.AddressToHex(*tx.Data.Recipient)//Recipient收件人
	//*tx.Data与tx.Data有啥区别
	//以下两个对比
	fromAddr := types.AddressToHex(craft.Address(from_))
	txFromAddr := types.AddressToHex(*tx.Data.From)
	//verify args
	//进行对比与传进来的值
	if amount != tx.Data.Amount.Uint64() || contractAddr != txToAddr || fromAddr != txFromAddr || to != targetToInput{
		//传入的值                              固定值     tx.Data.Recipient tx里面得到的sender
		//后者是*tx.Data里面的

		//*********这里是在比传入的值与什么的值？？？？？？！！！！！！！
		return errors.New("tx args not matched tx's"), 0
	}
	//如果有一个不匹配就报错
	log.Info("ReceiveFunds_verify_success, targetToAddr=%s, amount=%d, srcChainId=%d", to, amount, srcChainId)
	return nil, 1
}

// Register register a rpc route
func Register(methodName string, f *RPCFunc) error {
	paramStr := ""
	for _, arg := range f.args {
		switch arg.Kind() {
		case reflect.Uint64:
			paramStr += "uint64,"
		case reflect.String:
			paramStr += "string,"
		case reflect.Slice:
			if reflect.Uint8 != arg.Elem().Kind() {
				return errors.New("unsupported arg type")
			}
		}
	}
	if len(paramStr) > 0 {
		paramStr = paramStr[:len(paramStr)-1]
	}
	methodHash := util.Hash([]byte(methodName + "(" + paramStr + ")"))[:4]
	routes[string(methodHash)] = f
	return nil
}
//解析方法名
func Handler(input []byte) ([]byte, error) {
	method := util.ExtractMethodHash(input)//根据hash值找到方法名，通过提取input的前四个字节，即方法hash代号
	rpcFunc := routes[string(method)]//根据找到的方法创建了一个方法
	if rpcFunc == nil {
		return nil, errors.New("routes not found")
	}

	args, err := inputParamsToArgs(rpcFunc, input[len(method):])
	//这一步只是得到参数，并没有运行
	if err != nil {
		return nil, err
	}

	log.Info("contract RPC method: %s", util.BytesToHex(method))
	returns := rpcFunc.f.Call(args)    //把入参放入找到的方法并执行这个函数
	return encodeResult(returns)
}

// Covert an http query to a list of properly typed values.
// To be properly decoded the arg must be a concrete type from tendermint (if its an interface).
func inputParamsToArgs(rpcFunc *RPCFunc, input []byte) ([]reflect.Value, error) {
	args := make([]interface{}, 0)
	for _, argT := range rpcFunc.args {
		args = append(args, reflect.New(argT).Interface())
	}
	err := util.ExtractParam(input, args...)
	if err != nil {
		return nil, err
	}

	argVs := make([]reflect.Value, 0)
	for _, arg := range args {
		argVs = append(argVs, reflect.ValueOf(arg).Elem())
	}
	return argVs, nil
}

// NOTE: assume returns is result struct and error. If error is not nil, return it
func encodeResult(returns []reflect.Value) ([]byte, error) {
	errV := returns[0]
	if errV.Interface() != nil {
		return nil, errors.New(fmt.Sprintf("%v", errV.Interface()))
	}
	returns = returns[1:]
	rvs := make([]interface{}, 0)
	for _, rv := range returns {
		// the result is a registered interface,
		// we need a pointer to it so we can marshal with type byte
		rvp := reflect.New(rv.Type())
		rvp.Elem().Set(rv)
		rvs = append(rvs, rvp.Elem().Interface())
	}
	return util.EncodeReturnValue(rvs...)
}

// RPCFunc contains the introspected type information for a function
type RPCFunc struct {
	f       reflect.Value  // underlying rpc function
	args    []reflect.Type // type of each function arg
	returns []reflect.Type // type of each return arg
}

// NewRPCFunc create a new RPCFunc instance
func NewRPCFunc(f interface{}) *RPCFunc {
	return &RPCFunc{
		f:       reflect.ValueOf(f),
		args:    funcArgTypes(f),
		returns: funcReturnTypes(f),
	}
}

// return a function's argument types
func funcArgTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumIn()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.In(i)
	}
	return typez
}

// return a function's return types
func funcReturnTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumOut()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.Out(i)
	}
	return typez
}

//----------------------------------------------------------------------------------------------------------
//自己加的其他方法
func excBlockLimit( limit uint64 ,chainFlag string)(uint64){

	string_slice := strings.Split(chainFlag, ":")
	var ipData string =string_slice[0]
	var portData string =string_slice[1]
	web, _ := wutils.NewWeb3(ipData, portData, false)

	//查询过块高
	targetChainBloCKNum_, err:= web.Eth.BlockNumber()
	if(err!=nil){
		// 0有错，1没有，2过了
		return  0
	}
	targetChainBloCKNum := targetChainBloCKNum_.Uint64()
	if(targetChainBloCKNum>=limit){
		return  2

	}
	return  1

}



//out2的底层方法
func out2(txHash string, limit uint64,chainFlag string)(error, uint64){

	//发送http请求得到数据---
	url:="http://"+chainFlag+"/eth_getTransactionReceipt?hash=%22"+txHash+"%22"
	//http://192.168.160.129:47768/eth_getTransactionReceipt?hash=%220xe0ae85da4848d8725ac43ba0ad95cc5cfe18d1c80773ebd85239cd7bf91898ef%22

	resp, err:=http.Get(url)
	if err != nil {
		return err, 1
	}
	defer resp.Body.Close()
	body, err:= ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, 1
	}

	var recipient RecipientData
	err = json.Unmarshal(body, &recipient)
	if  err != nil {
		return err, 1
	}
	//有交易是"0x1",没有是""
	txState:=recipient.Result.Status
	blockLimitState1:=excBlockLimit(limit,chainFlag)

	//todo交易未上链
	if(txState==""){
		// 0有错，1没有，2过了

		if(blockLimitState1==0){
			return err, 1
		}else if(blockLimitState1==1){
			return err, 1
		}
		return err, 3
	}

	//没有logs就是nil
	logs_:=recipient.Result.Logs

	if(logs_==nil){

		return err, 3

	}

	state:=string(logs_[0].Data)

	//A 0,E 1,I 2,M 3
	value1:="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE="
	//value2:="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAI="
	//value3:="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAM="

	if(state==value1) {
		return err, 2
	}
	return err, 3
	//3交易失败  2交易成功 1交易状态未知
}

//in4的底层查询方法
func In4(txHash string, limit uint64,chainFlag string)(error, uint64){

	//发送http请求得到数据---
	url:="http://"+chainFlag+"/eth_getTransactionReceipt?hash=%22"+txHash+"%22"
	//http://192.168.160.129:47768/eth_getTransactionReceipt?hash=%220xe0ae85da4848d8725ac43ba0ad95cc5cfe18d1c80773ebd85239cd7bf91898ef%22

	resp, err:=http.Get(url)
	if err != nil {
		return err, 1
	}
	defer resp.Body.Close()
	body, err:= ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, 1
	}

	var recipient RecipientData
	err = json.Unmarshal(body, &recipient)
	if  err != nil {
		return err, 1
	}
	//有交易是"0x1",没有是""
	txState:=recipient.Result.Status
	blockLimitState1:=excBlockLimit(limit,chainFlag)

	//todo交易未上链
	if(txState==""){
		// 0有错，1没有，2过了

		if(blockLimitState1==0){
			return err, 1
		}else if(blockLimitState1==1){
			return err, 1
		}
		return err, 3
	}

	//没有logs就是nil
	logs_:=recipient.Result.Logs

	if(logs_==nil){

		return err, 3

	}

	state:=string(logs_[0].Data)

	//A 0,E 1,I 2,M 3
	value1:="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE="
	value2:="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAI="
	//value3:="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAM="


	//1 不明白  2成功  3失败
	if(state==value2) {
		return err, 2
	}else if(state==value1){
		return err, 1
	}
	return err, 3

	//3交易失败  2交易成功 1交易状态未知


}
//只传入交易hash与 ip就可以用
func generate(txHash string,chainFlag string)(error, uint64){
	//发送http请求得到数据---
	url:="http://"+chainFlag+"/eth_getTransactionReceipt?hash=%22"+txHash+"%22"
	//http://192.168.160.129:47768/eth_getTransactionReceipt?hash=%220xe0ae85da4848d8725ac43ba0ad95cc5cfe18d1c80773ebd85239cd7bf91898ef%22

	resp, err:=http.Get(url)
	if err != nil {
		return err, 1
	}
	defer resp.Body.Close()
	body, err:= ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, 1
	}

	var recipient RecipientData
	err = json.Unmarshal(body, &recipient)
	if  err != nil {
		return err, 1
	}
	//有交易是"0x1",没有是""
	txState:=recipient.Result.Status

	//todo交易未上链
	if(txState==""){
		return err, 0
	}
	return err, 1
	//1交友上链  2交易未上链


}




