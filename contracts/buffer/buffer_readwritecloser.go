package buffer

import (
	"github.com/ethereum/go-ethereum/common"
	"io"
)

type SystemBufferReadWriterCloser struct {
	sysBufferContract *SystemBufferContract
	cursor            uint64
}

//NewSystemBufferReadWriterCloser create a new instance
func NewSystemBufferReadWriterCloser(sysBufferContract *SystemBufferContract) *SystemBufferReadWriterCloser {
	return &SystemBufferReadWriterCloser{
		cursor:            0,
		sysBufferContract: sysBufferContract,
	}
}

func (this *SystemBufferReadWriterCloser) Read(data []byte) (n int, err error) {
	totalLen := this.sysBufferContract.Length()
	if this.cursor >= totalLen {
		return 0, io.EOF
	}
	size := uint64(0)
	if totalLen-this.cursor > uint64(len(data)) {
		size = uint64(len(data))
	} else {
		size = totalLen - this.cursor
	}
	ret, err := this.sysBufferContract.Read(this.cursor, size)
	n = copy(data, ret)
	this.cursor += uint64(n)
	return n, nil
}

func (this *SystemBufferReadWriterCloser) Write(data []byte) (n int, err error) {
	len, err := this.sysBufferContract.Write(data)
	return int(len), err
}

func (this *SystemBufferReadWriterCloser) Close() error {
	return this.sysBufferContract.Close()
}

func (this *SystemBufferReadWriterCloser) ContractAddress() common.Address {
	return this.sysBufferContract.Address()
}
