// Copyright (c) 2016, Alan Chen
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package web3

import (
	"encoding/json"
	"math/big"
	"github.com/ethereum/go-ethereum/web3go/common"
)

type jsonBlock struct {
	Number          json.Number    `json:"number"`
	Hash            common.Hash    `json:"hash"`
	ParentHash      common.Hash    `json:"parentHash"`
	Nonce           common.Hash    `json:"nonce"`
	Sha3Uncles      common.Hash    `json:"sha3Uncles"`
	Bloom           common.Hash    `json:"logsBloom"`
	TransactionRoot common.Hash    `json:"transactionsRoot"`
	StateRoot       common.Hash    `json:"stateRoot"`
	Miner           common.Address `json:"miner"`
	Difficulty      json.Number    `json:"difficulty"`
	TotalDifficulty json.Number    `json:"totalDifficulty"`
	ExtraData       common.Hash    `json:"extraData"`
	Size            json.Number    `json:"size"`
	GasLimit        json.Number    `json:"gasLimit"`
	GasUsed         json.Number    `json:"gasUsed"`
	Timestamp       json.Number    `json:"timestamp"`
	Transactions    []common.Hash  `json:"transactions"`
	Uncles          []common.Hash  `json:"uncles"`
}

func (b *jsonBlock) ToBlock() (block *common.Block) {
	block = &common.Block{}
	block.Number = jsonNumbertoInt(b.Number)
	block.Hash = b.Hash
	block.ParentHash = b.ParentHash
	block.Nonce = b.Nonce
	block.Sha3Uncles = b.Sha3Uncles
	block.Bloom = b.Bloom
	block.TransactionRoot = b.TransactionRoot
	block.StateRoot = b.StateRoot
	block.Miner = b.Miner
	block.Difficulty = jsonNumbertoInt(b.Difficulty)
	block.TotalDifficulty = jsonNumbertoInt(b.TotalDifficulty)
	block.ExtraData = b.ExtraData
	block.Size = jsonNumbertoInt(b.Size)
	block.GasLimit = jsonNumbertoInt(b.GasLimit)
	block.GasUsed = jsonNumbertoInt(b.GasUsed)
	block.Timestamp = jsonNumbertoInt(b.Timestamp)
	block.Transactions = b.Transactions
	block.Uncles = b.Uncles
	return block
}

type jsonTransaction struct {
	Hash             string    `json:"hash"`
	Nonce            string    `json:"nonce"`
	BlockHash        string    `json:"blockHash"`
	BlockNumber      string    `json:"blockNumber"`
	TransactionIndex string    `json:"transactionIndex"`
	From             string    `json:"from"`
	To               string    `json:"to"`
	Gas              string    `json:"gas"`
	GasPrice         string    `json:"gasPrice"`
	Value            string    `json:"value"`
	Data             string         `json:"input"`
	R                string         `json:"r"`
	S                string         `json:"s"`
	V                string         `json:"v"`
}

func (t *jsonTransaction) ToTransaction() (tx *common.Transaction) {
	tx = &common.Transaction{}

	tx.Hash = common.StringToHash(t.Hash)
	tx.Nonce = jsonNumbertoInt(json.Number(t.Nonce))
	tx.BlockHash = common.StringToHash(t.BlockHash)
	tx.BlockNumber = jsonNumbertoInt(json.Number(t.BlockNumber))
	txIndex, _ := json.Number(t.TransactionIndex).Int64()
	tx.TransactionIndex = uint64(txIndex)
	tx.From = common.StringToAddress(t.From)
	tx.To = common.StringToAddress(t.To)
	tx.Gas = jsonNumbertoInt(json.Number(t.Gas))
	tx.GasPrice = jsonNumbertoInt(json.Number(t.GasPrice))
	tx.Value = jsonNumbertoInt(json.Number(t.Value))
	tx.Data = common.HexToBytes(t.Data)
	tx.R = common.HexToBytes(t.R)
	tx.S = common.HexToBytes(t.S)
	tx.V = common.HexToBytes(t.V)

	//tx.Hash = t.Hash
	//tx.Nonce = t.Nonce
	//tx.BlockHash = t.BlockHash
	//tx.BlockNumber = jsonNumbertoInt(t.BlockNumber)
	//tx.TransactionIndex = t.TransactionIndex
	//tx.From = t.From
	//tx.To = t.To
	//tx.Gas = jsonNumbertoInt(t.Gas)
	//tx.GasPrice = jsonNumbertoInt(t.GasPrice)
	//tx.Value = jsonNumbertoInt(t.Value)
	//tx.Data = t.Data
	//tx.R = t.R
	//tx.S = t.S
	//tx.V = t.V

	return tx
}

func (tx *jsonTransaction) Unmarshal(data []byte) error{
	var result map[string]string
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	tx.Hash = result["hash"]
	tx.Nonce = result["nonce"]
	tx.BlockHash = result["blockHash"]
	tx.BlockNumber = result["blockNumber"]
	tx.TransactionIndex = result["transactionIndex"]
	tx.From = result["from"]
	tx.To = result["to"]
	tx.Gas = result["gas"]
	tx.GasPrice = result["gasPrice"]
	tx.Value = result["value"]
	tx.Data = result["input"]
	tx.R = result["r"]
	tx.S = result["s"]
	tx.V = result["v"]

	return nil
}

func (tx *jsonTransaction) String() string {
	jsonBytes, _ := json.Marshal(tx)
	return string(jsonBytes)
}

type jsonTransactionReceipt struct {
	BlockHash         common.Hash     	`json:"blockHash"`
	BlockNumber       json.Number 		`json:"blockNumber"`
	Hash              common.Hash    	`json:"transactionHash"`
	TransactionIndex  uint64   			`json:"transactionIndex"`
	From			  common.Address 	`json:"from"`
	To				  common.Address 	`json:"to"`
	Root			  []byte  			`json:"root"`
	Status			  json.Number 		`json:"status"`
	GasUsed           json.Number 		`json:"gasUsed"`
	CumulativeGasUsed json.Number 		`json:"cumulativeGasUsed"`
	LogsBloom		  []byte   			`json:"logsBloom"`
	Logs              []jsonLog    		`json:"logs"`
	ContractAddress   common.Address  	`json:"contractAddress"`
}

func (r *jsonTransactionReceipt) UnmarshalJSON(input []byte) error {
	type jsonTransactionReceipt struct{
		BlockHash         string     	`json:"blockHash"`
		BlockNumber       json.Number 	`json:"blockNumber"`
		Hash              string    	`json:"transactionHash"`
		TransactionIndex  json.Number   `json:"transactionIndex"`
		From			  string 		`json:"from"`
		To				  string 		`json:"to"`
		Root			  string  		`json:"root"`
		Status			  json.Number 	`json:"status"`
		GasUsed           json.Number 	`json:"gasUsed"`
		CumulativeGasUsed json.Number 	`json:"cumulativeGasUsed"`
		LogsBloom		  string   		`json:"logsBloom"`
		Logs              []jsonLog    	`json:"logs"`
		ContractAddress   string  	    `json:"contractAddress"`
	}
	var dec jsonTransactionReceipt
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	r.BlockHash = common.StringToHash(dec.BlockHash)
	r.BlockNumber = dec.BlockNumber
	r.Hash = common.StringToHash(dec.Hash)
	index, _ := dec.TransactionIndex.Int64()
	r.TransactionIndex = uint64(index)
	r.From = common.StringToAddress(dec.From)
	r.To = common.StringToAddress(dec.To)
	r.Root = common.HexToBytes(dec.Root)
	r.Status = dec.Status
	r.GasUsed = dec.GasUsed
	r.CumulativeGasUsed = dec.CumulativeGasUsed
	r.LogsBloom = common.HexToBytes(dec.LogsBloom)
	r.Logs = dec.Logs
	r.ContractAddress = common.StringToAddress(dec.ContractAddress)

	return nil
}

func (r *jsonTransactionReceipt) ToTransactionReceipt() (receipt *common.TransactionReceipt) {
	receipt = &common.TransactionReceipt{}
	receipt.Hash = r.Hash
	receipt.TransactionIndex = r.TransactionIndex
	receipt.BlockNumber = jsonNumbertoInt(r.BlockNumber)
	receipt.BlockHash = r.BlockHash
	receipt.CumulativeGasUsed = jsonNumbertoInt(r.CumulativeGasUsed)
	receipt.GasUsed = jsonNumbertoInt(r.GasUsed)
	receipt.ContractAddress = r.ContractAddress
	receipt.Status = jsonNumbertoInt(r.Status)
	receipt.Logs = make([]common.Log, 0)
	for _, l := range r.Logs {
		receipt.Logs = append(receipt.Logs, l.ToLog())
	}
	return receipt
}

type jsonLog struct {
	LogIndex         uint64         `json:"logIndex"`
	BlockNumber      json.Number    `json:"blockNumber"`
	BlockHash        common.Hash    `json:"blockHash"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	TransactionIndex uint64         `json:"transactionIndex"`
	Address          common.Address `json:"address"`
	Data             []byte         `json:"data"`
	Topics           common.Topics  `json:"topics"`
}

func (l jsonLog) ToLog() (log common.Log) {
	log = common.Log{}
	log.LogIndex = l.LogIndex
	log.BlockNumber = jsonNumbertoInt(l.BlockNumber)
	log.BlockHash = l.BlockHash
	log.TransactionHash = l.TransactionHash
	log.TransactionIndex = l.TransactionIndex
	log.Address = l.Address
	log.Data = l.Data
	log.Topics = l.Topics
	return log
}

func jsonNumbertoInt(data json.Number) *big.Int {
	f := big.NewFloat(0.0)
	f.SetString(string(data))
	result, _ := f.Int(nil)
	return result
}
