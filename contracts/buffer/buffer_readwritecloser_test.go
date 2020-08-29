package buffer
//
//import (
//	"github.com/DSiSc/craft/types"
//	"github.com/DSiSc/evm-NG/common/hexutil"
//	"github.com/DSiSc/monkey"
//	"github.com/stretchr/testify/assert"
//	"reflect"
//	"testing"
//)
//
//var mockData, _ = hexutil.Decode("0x11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
//	"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
//	"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
//	"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
//	"111111111111111111111111116666")
//
//func mockSystemBuffer() *SystemBufferContract {
//	sysBuffer := &SystemBufferContract{}
//	return sysBuffer
//}
//
//func TestNewSystemBufferReadWriterCloser(t *testing.T) {
//	assert := assert.New(t)
//	sysRWC := NewSystemBufferReadWriterCloser(mockSystemBuffer())
//	assert.NotNil(sysRWC)
//}
//
//func TestSystemBufferReadWriterCloser_Read(t *testing.T) {
//	defer monkey.UnpatchAll()
//	assert := assert.New(t)
//	sysRWC := NewSystemBufferReadWriterCloser(mockSystemBuffer())
//	assert.NotNil(sysRWC)
//
//	monkey.PatchInstanceMethod(reflect.TypeOf(sysRWC.sysBufferContract), "Length", func(sysBuffer *SystemBufferContract) uint64 {
//		return uint64(len(mockData))
//	})
//	monkey.PatchInstanceMethod(reflect.TypeOf(sysRWC.sysBufferContract), "Read", func(sysBuffer *SystemBufferContract, offset, size uint64) ([]byte, error) {
//		return mockData[int(offset):int(offset+size)], nil
//	})
//
//	data := make([]byte, 10)
//	n1, err := sysRWC.Read(data)
//	assert.Nil(err)
//	assert.Equal(len(data), n1)
//
//	data = make([]byte, len(mockData))
//	n2, err := sysRWC.Read(data)
//	assert.Nil(err)
//	assert.Equal(len(mockData)-n1, n2)
//
//	data = make([]byte, 10)
//	_, err = sysRWC.Read(data)
//	assert.NotNil(err)
//}
//
//func TestSystemBufferReadWriterCloser_Write(t *testing.T) {
//	defer monkey.UnpatchAll()
//	assert := assert.New(t)
//	sysRWC := NewSystemBufferReadWriterCloser(mockSystemBuffer())
//	assert.NotNil(sysRWC)
//
//	monkey.PatchInstanceMethod(reflect.TypeOf(sysRWC.sysBufferContract), "Write", func(sysBuffer *SystemBufferContract, data []byte) (uint64, error) {
//		return uint64(len(data)), nil
//	})
//
//	data := []byte("Hello, World")
//	n, err := sysRWC.Write(data)
//	assert.Nil(err)
//	assert.Equal(len(data), n)
//}
//
//func TestSystemBufferContract_Close(t *testing.T) {
//	defer monkey.UnpatchAll()
//	assert := assert.New(t)
//	sysRWC := NewSystemBufferReadWriterCloser(mockSystemBuffer())
//	assert.NotNil(sysRWC)
//
//	monkey.PatchInstanceMethod(reflect.TypeOf(sysRWC.sysBufferContract), "Close", func(sysBuffer *SystemBufferContract) error {
//		return nil
//	})
//	assert.Nil(sysRWC.Close())
//}
//
//func TestSystemBufferReadWriterCloser_ContractAddress(t *testing.T) {
//	defer monkey.UnpatchAll()
//	assert := assert.New(t)
//	sysRWC := NewSystemBufferReadWriterCloser(mockSystemBuffer())
//	assert.NotNil(sysRWC)
//
//	monkey.PatchInstanceMethod(reflect.TypeOf(sysRWC.sysBufferContract), "Address", func(sysBuffer *SystemBufferContract) types.Address {
//		return SystemBufferAddr
//	})
//	assert.Equal(SystemBufferAddr, sysRWC.ContractAddress())
//}
