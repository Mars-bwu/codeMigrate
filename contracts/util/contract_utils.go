package util

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/contracts/constant"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"math/big"
	"reflect"
)

const HashLenght = 32

var (
	UnSupportedTypeError  = errors.New("unsupported arg type")
	InvalidUnmarshalError = errors.New("invalid unmarshal error")
)

// ExtractMethodHash extract method hash from input
func ExtractMethodHash(input []byte) []byte {
	return input[:4]
}

// Hash return the hash of the data
func Hash(data []byte) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// ExtractParam extract string params from input
func ExtractParam(input []byte, args ...interface{}) error {
	for i := 0; i < len(args); i++ {
		rv := reflect.ValueOf(args[i])
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return InvalidUnmarshalError
		}

		switch rv.Elem().Type().Kind() {
		case reflect.String:
			arg := string(extractDynamicTypeData(input, i))
			rv.Elem().SetString(arg)
		case reflect.Slice:
			if reflect.Uint8 != rv.Elem().Type().Elem().Kind() {
				return UnSupportedTypeError
			}
			arg := extractDynamicTypeData(input, i)
			rv.Elem().SetBytes(arg)
		case reflect.Uint64:
			arg, _ := math.ParseUint64(hexutil.Encode(input[i*constant.EvmWordSize : (i+1)*constant.EvmWordSize]))
			rv.Elem().SetUint(arg)
		case reflect.Array:
			arg := arrayByte20(input, i)
			rv.Elem().SetBytes(arg)
		default:
			return UnSupportedTypeError
		}
	}
	return nil
}

// extract dynamic type data
func extractDynamicTypeData(totalInput []byte, varIndex int) []byte {
	offset, _ := math.ParseUint64(hexutil.Encode(totalInput[varIndex*constant.EvmWordSize : (varIndex+1)*constant.EvmWordSize]))
	if offset >= uint64(len(totalInput)) {
		// address type
		addr := types.BytesToAddress(arrayByte20(totalInput, varIndex))
		addrStr := types.AddressToHex(addr)
		return []byte(addrStr)
	}
	dataLen, _ := math.ParseUint64(hexutil.Encode(totalInput[offset : offset+constant.EvmWordSize]))
	argStart := offset + constant.EvmWordSize
	argEnd := argStart + dataLen
	return totalInput[argStart:argEnd]
}

// EncodeReturnValue encode the return value to the format needed by evm将返回值编码为EVM返回值需要的格式
func EncodeReturnValue(retVals ...interface{}) ([]byte, error) {
	retPre := make([]byte, 0)
	retData := make([]byte, 0)
	preOffsetPadding := len(retVals) * constant.EvmWordSize
	for _, retVal := range retVals {
		retType := reflect.TypeOf(retVal)
		switch retType.Kind() {
		case reflect.String:
			offset := preOffsetPadding + len(retData)
			retPre = append(retPre, math.PaddedBigBytes(big.NewInt(int64(offset)), constant.EvmWordSize)...)
			retData = append(retData, encodeString(retVal.(string))...)
		case reflect.Slice:
			if reflect.Uint8 != retType.Elem().Kind() {
				return nil, UnSupportedTypeError
			}
			offset := preOffsetPadding + len(retData)
			retPre = append(retPre, math.PaddedBigBytes(big.NewInt(int64(offset)), constant.EvmWordSize)...)
			retData = append(retData, encodeBytes(retVal.([]byte))...)
		case reflect.Uint64:
			retPre = append(retPre, math.PaddedBigBytes(big.NewInt(0).SetUint64(retVal.(uint64)), constant.EvmWordSize)...)
		case reflect.Array:
			if retType.AssignableTo(reflect.TypeOf(types.Address{})) {
				addr := retVal.(types.Address)
				retPre = append(retPre, common.LeftPadBytes(addr[:], constant.EvmWordSize)...)
			}
		default:
			return nil, errors.New("unsupported return type")
		}
	}
	return append(retPre, retData...), nil
}

// encode the string to the format needed by evm
func encodeString(val string) []byte {
	return encodeBytes([]byte(val))
}

// encode the byte array to the format needed by evm
func encodeBytes(val []byte) []byte {
	ret := make([]byte, 0)
	ret = append(ret, math.PaddedBigBytes(big.NewInt(int64(len(val))), constant.EvmWordSize)...)
	for i := 0; i < len(val); {
		if (len(val) - i) > constant.EvmWordSize {
			ret = append(ret, val[i:i+constant.EvmWordSize]...)
			i += constant.EvmWordSize
		} else {
			ret = append(ret, common.RightPadBytes(val[i:], constant.EvmWordSize)...)
			i += len(val)
		}
	}
	return ret
}

func arrayByte20(totalInput []byte, varIndex int) []byte {
	offset := varIndex*constant.EvmWordSize + constant.AddressOffset
	addrByte := totalInput[offset : (varIndex+1)*constant.EvmWordSize]
	return []byte(addrByte)
}

