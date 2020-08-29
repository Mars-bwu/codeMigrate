package vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/buffer"
	"github.com/ethereum/go-ethereum/contracts/rpc"

	//	"github.com/DSiSc/evm-NG/system/contract/buffer"
	//	"github.com/DSiSc/evm-NG/system/contract/rpc"
	//	"github.com/DSiSc/evm-NG/system/contract/storage"
)
//
// SysContractExecutionFunc system contract execute function
//合约方法类型
//定义函数类型的固定写法
type SysContractExecutionFunc func(interpreter *EVM, contract ContractRef, input []byte) ([]byte, error)
//https://blog.csdn.net/tzs919/article/details/53571632函数类型
// system call routes
//                     地址类型 =>      映射到函数类型
var routes = make(map[common.Address]SysContractExecutionFunc)
//make用于创建映射

//这里初始化
func init() {
	//与映射一样这个地址关联了一个函数类型
	//没看懂但是不管怎么样，这里能得到一个方法？？？？
	routes[buffer.SystemBufferAddr] = func(execEvm *EVM, contract ContractRef, input []byte) ([]byte, error) {
		systemBuffer := buffer.NewSystemBufferContract(execEvm.StateDB)
		return buffer.BufferExecute(systemBuffer, input)
	}

	//routes[storage.TencentCosAddr] = func(execEvm *EVM, caller ContractRef, input []byte) ([]byte, error) {
	//	systemBuffer := buffer.NewSystemBufferContract(execEvm.StateDB)
	//	systemBufferReadWriter := buffer.NewSystemBufferReadWriterCloser(systemBuffer)
	//	tencentCos := storage.NewTencentCosContract(systemBufferReadWriter)
	//	return storage.CosExecute(tencentCos, input)
	//}
	//根据合约去找具体的方法
	routes[rpc.RpcContractAddr] = func(execEvm *EVM, caller ContractRef, input []byte) ([]byte, error) {
		return rpc.Handler(input)
		//rpc.Handler能找到方法名
		//这个arg----input是输入的参数----这里是解析出他的方法名
	}
}

//IsSystemContract check the contract with specified address is system contract
func IsSystemContract(addr common.Address) bool {
	return routes[addr] != nil
}

// GetSystemContractExecFunc get system contract execution function by address
//通过智能合约地址获得合约方法
func GetSystemContractExecFunc(addr common.Address) SysContractExecutionFunc {
	//                                                可见返回值是一个方法
	return routes[addr]
}
