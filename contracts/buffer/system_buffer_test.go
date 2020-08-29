package buffer

import (
	"bytes"
	"encoding/binary"
	"github.com/DSiSc/evm-NG/common/hexutil"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/repository"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewSystemBufferContract(t *testing.T) {
	assert := assert.New(t)
	db := &repository.Repository{}
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)

}

func TestBufferExecute(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	db := &repository.Repository{}
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)

	data := []byte{0x11, 0x11, 0x11}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		if bytes.Equal([]byte(systemBufferCacheKey), key) {
			val := make([]byte, 8)
			binary.BigEndian.PutUint64(val, 3)
			return val, nil
		}
		return data, nil
	})
	input, _ := hexutil.Decode("0xae0bf88300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003")
	ret, err := BufferExecute(bc, input)
	assert.Nil(err)
	expectRet, _ := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000031111110000000000000000000000000000000000000000000000000000000000")
	assert.Equal(expectRet, ret)
}

func TestBufferExecute1(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	db := &repository.Repository{}
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)

	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return nil, nil
	})
	input, _ := hexutil.Decode("0xc35789cc")
	_, err := BufferExecute(bc, input)
	assert.Nil(err)

	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		val := make([]byte, 8)
		binary.BigEndian.PutUint64(val, uint64(1))
		return val, nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Delete", func(chain *repository.Repository, key []byte) error {
		return nil
	})
	_, err = BufferExecute(bc, input)
	assert.Nil(err)
}

func TestBufferExecute2(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	db := &repository.Repository{}
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	input, _ := hexutil.Decode("0x82172882")
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return nil, nil
	})
	expectRet, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000000")
	ret, err := BufferExecute(bc, input)
	assert.Nil(err)
	assert.Equal(expectRet, ret)

	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		val := make([]byte, 8)
		binary.BigEndian.PutUint64(val, uint64(1))
		return val, nil
	})
	expectRet, _ = hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000001")
	ret, err = BufferExecute(bc, input)
	assert.Nil(err)
	assert.Equal(expectRet, ret)
}

func TestBufferExecute3(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	db := &repository.Repository{}
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)

	var data []byte
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		if !bytes.Equal([]byte(systemBufferCacheKey), key) {
			data = value
		}
		return nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return nil, nil
	})
	input, _ := hexutil.Decode("0x5f10585d000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000031111110000000000000000000000000000000000000000000000000000000000")
	_, err := BufferExecute(bc, input)
	assert.Nil(err)
	expectData := []byte{0x11, 0x11, 0x11}
	assert.Equal(expectData, data)
}

func TestSystemBufferContract_Read(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	_, err := bc.Read(0, 1)
	assert.NotNil(err)

	data := []byte{0x11, 0x11, 0x11}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		if bytes.Equal([]byte(systemBufferCacheKey), key) {
			val := make([]byte, 8)
			binary.BigEndian.PutUint64(val, 3)
			return val, nil
		}
		return data, nil
	})
	ret, err := bc.Read(0, 3)
	assert.Nil(err)
	assert.NotNil(data, ret)
}

func TestSystemBufferContract_Write(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	data, _ := hexutil.Decode("0x111111111111111111111111111111111111111111116666")
	len, err := bc.Write(data)
	assert.Nil(err)
	assert.Equal(uint64(0x18), len)
}

func TestSystemBufferContract_Read1(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	data, err := bc.Read(0, 1)
	assert.NotNil(err)

	data, _ = hexutil.Decode("0x1111111111111111116666")
	lenD, err := bc.Write(data)
	assert.Nil(err)
	assert.Equal(uint64(11), lenD)

	data, err = bc.Read(9, 2)
	assert.Nil(err)
	assert.Equal([]byte{0x66, 0x66}, data)
}

func TestSystemBufferContract_Read2(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	data, err := bc.Read(0, 1)
	assert.NotNil(err)

	data, _ = hexutil.Decode("0x11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111116666")
	lenD, err := bc.Write(data)
	assert.Nil(err)
	assert.Equal(uint64(259), lenD)

	data, err = bc.Read(257, 2)
	assert.Nil(err)
	assert.Equal([]byte{0x66, 0x66}, data)
}

func TestSystemBufferContract_Write1(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	data, _ := hexutil.Decode("0x111111111111111111111111111111111111111111116666")
	len, err := bc.Write(data)
	assert.Nil(err)
	assert.Equal(uint64(0x18), len)
}

func TestSystemBufferContract_Write2(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	data, _ := hexutil.Decode("0x111111111111111111111111111111111111111111116666")
	len, err := bc.Write(data)
	assert.Nil(err)
	assert.Equal(uint64(0x18), len)
	assert.Equal(uint64(0x18), binary.BigEndian.Uint64(lowLevelCache[systemBufferCacheKey]))

	data, _ = hexutil.Decode("0x1234")
	len, err = bc.Write(data)
	assert.Nil(err)
	assert.Equal(uint64(2), len)
	assert.Equal(uint64(0x1A), binary.BigEndian.Uint64(lowLevelCache[systemBufferCacheKey]))
}

func TestSystemBufferContract_Length(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)

	saveLen := bc.Length()
	assert.Equal(uint64(0), saveLen)

	data, _ := hexutil.Decode("0x111111111111111111111111111111111111111111116666")
	len, err := bc.Write(data)
	assert.Nil(err)
	assert.Equal(uint64(0x18), len)
	saveLen = bc.Length()
	assert.Equal(uint64(0x18), saveLen)

	data, _ = hexutil.Decode("0x1234")
	len, err = bc.Write(data)
	assert.Nil(err)
	saveLen = bc.Length()
	assert.Equal(uint64(0x1A), saveLen)
}

func TestSystemBufferContract_Delete(t *testing.T) {
	defer monkey.UnpatchAll()
	lowLevelCache := make(map[string][]byte)
	assert := assert.New(t)
	db := &repository.Repository{}
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Get", func(chain *repository.Repository, key []byte) ([]byte, error) {
		return lowLevelCache[string(key)], nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Put", func(chain *repository.Repository, key []byte, value []byte) error {
		lowLevelCache[string(key)] = value
		return nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(db), "Delete", func(chain *repository.Repository, key []byte) error {
		delete(lowLevelCache, string(key))
		return nil
	})
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)

	data, _ := hexutil.Decode("0x111111111111111111111111111111111111111111116666")
	_, err := bc.Write(data)
	assert.Nil(err)

	err = bc.Close()
	assert.Nil(err)
	assert.Equal(0, len(lowLevelCache))

	data, _ = hexutil.Decode("0x11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111116666")
	_, err = bc.Write(data)
	assert.Nil(err)

	err = bc.Close()
	assert.Nil(err)
	assert.Equal(0, len(lowLevelCache))
}

func TestSystemBufferContract_Address(t *testing.T) {
	assert := assert.New(t)
	db := &repository.Repository{}
	bc := NewSystemBufferContract(db)
	assert.NotNil(bc)
	assert.Equal(SystemBufferAddr, bc.Address())
}
