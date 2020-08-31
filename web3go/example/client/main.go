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

package main

import (
	"flag"
	"fmt"
	"github.com/DSiSc/web3go/common"
	"github.com/DSiSc/web3go/provider"
	"github.com/DSiSc/web3go/rpc"
	"github.com/DSiSc/web3go/web3"
	"math/big"
)

var hostname = flag.String("hostname", "127.0.0.1", "The ethereum client RPC host")
var port = flag.String("port", "47768", "The ethereum client RPC port")
var verbose = flag.Bool("verbose", false, "Print verbose messages")

func main() {
	flag.Parse()

	if *verbose {
		fmt.Printf("Connect to %s:%s\n", *hostname, *port)
	}

	provider := provider.NewHTTPProvider(*hostname+":"+*port, rpc.GetDefaultMethod())
	web3 := web3.NewWeb3(provider)

	if accounts, err := web3.Eth.Accounts(); err == nil {
		for _, account := range accounts {
			fmt.Printf("%s\n", account.String())
		}
	} else {
		fmt.Printf("%v", err)
	}

	//from := common.NewAddress(common.HexToBytes("0x8025c2eeF50a15D29aC839Aed47c3c78F0cAC143"))
	to := common.NewAddress(common.HexToBytes("0x368AB89547Aad5604Fce277Cd6dB581851c337d5"))
    tx := common.NewTransactions(uint64(0), to, big.NewInt(100000), 100000, big.NewInt(0), make([]byte, 0))
    tx_encode, _ := tx.EncodeToRLP()
	res, _ := web3.Eth.SendRawTransaction(tx_encode)
	fmt.Println(res.String())
}
