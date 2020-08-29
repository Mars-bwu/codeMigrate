package util

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"

	"math/big"
)

// Lengths of hashes and addresses in bytes.
const (
	HashLength    = 32
	AddressLength = 20
)

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) common.Address {
	var a common.Address
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
	return a
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) common.Address { return BytesToAddress(FromHex(s)) }

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) common.Address { return BytesToAddress(b.Bytes()) }

// Big converts an address to a big integer.
func AddressToBig(addr common.Address) *big.Int { return new(big.Int).SetBytes(addr[:]) }

// AddressToBytes convert address to []byte
func AddressToBytes(addr common.Address) []byte {
	return addr[:]
}

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) common.Hash { return BytesToHash(FromHex(s)) }

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) common.Hash {
	var h common.Hash
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
	return h
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) common.Hash { return BytesToHash(b.Bytes()) }

// HashToBytes convert hash to []byte
func HashToBytes(hash common.Hash) []byte {
	return hash[:]
}

// HashToBig converts a hash to a big integer.
func HashToBig(hash common.Hash) *big.Int { return new(big.Int).SetBytes(hash[:]) }

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// Bytes2Hex returns the hexadecimal encoding of d.
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}
