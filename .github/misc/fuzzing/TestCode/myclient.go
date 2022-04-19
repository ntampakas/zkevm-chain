package main

import (
	"os"
	"context"
	"fmt"
	"zkevmchaintest/l1bridge"
	"log"
	"math/big"
	"math/rand"
	"time"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var passw0rd string
var TxHashes []common.Hash

func main() {

	doDeposits, _ := strconv.ParseBool(os.Args[1])
	doTxs, _ := strconv.ParseBool(os.Args[2])
	layer := os.Args[3]

	fmt.Printf("deposits: %T\nTxs: %T\nLayer: %T\n", doDeposits, doTxs, layer)

	// Random source for selecting account
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// keystore jsons location and accounts' passphrase
	ksdir := "./keystore"
	passw0rd := "password"
	//
	lOnebridge := "936a70c0b28532aa22240dce21f89a8399d6ac60"
	// new map with testnet urls
	tnurls := map[string]string{
		"layer1" : "https://zkevmchain-nick3.efprivacyscaling.org/rpc/l1",
		"layer2" : "https://zkevmchain-nick3.efprivacyscaling.org/rpc/l2",
	}

	// zChainTnUrls := map[string]string{
	// 	"layer1": "https://zkevmchaingeth1.efprivacyscaling.org/rpc/l1",
	// 	"layer2": "https://zkevmchaingeth1.efprivacyscaling.org/rpc/l2",
	// }

	//load keystore and initiate ethclients towards l1 and l2 networks
	_accounts, ks := LoadAccounts(ksdir)
	_ctx := context.Background()
	ethcl1, _ := ethclient.Dial(tnurls["layer1"])
	ethcl2, _ := ethclient.Dial(tnurls["layer2"])

	// zeros := make([]*big.Int, len(_accounts))

	//Create a new layer1 bridge instance
	bridgeAddress := common.HexToAddress(lOnebridge)
	bridge, err := l1bridge.NewZkEvmL1Bridge(bridgeAddress, ethcl1)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}

	// unlock keystore to enable Tx signing
	for _, account := range _accounts {
		ks.Unlock(account, passw0rd)
	}

	l1ChainID, err := ethcl1.NetworkID(_ctx)
	if err != nil {
		log.Fatal(err)
	}

	l2ChainID, err := ethcl2.NetworkID(_ctx)
	if err != nil {
		log.Fatal(err)
	}

	var zeros []*big.Int

	//Print l1 & l2 balances

	bal := GetBalances(_accounts, *ethcl1, *ethcl2, _ctx)

	for _, b := range bal {
		// fmt.Printf("account %x has %v funds in l1 and %v funds in l2\n", b.hexaddr, b.layer1Funds, b.layer2Funds)
		fmt.Printf("account %x has %v funds in l2 and %v funds in l1\n", b.hexaddr, b.layer2Funds, b.layer1Funds)
		if b.layer2Funds.Cmp(big.NewInt(0)) == 0 {
			zeros = append(zeros, b.layer2Funds)
			// fmt.Println("true")
		}
	}

	// for _, bbb := range bal {
	// 	fmt.Println(bbb.layer2Funds, " : ", reflect.TypeOf(bbb.layer2Funds), " : ", big.NewInt(0), " : ", reflect.TypeOf(big.NewInt(0)))
	// }

	// fmt.Println(zeros)

	//generate new Tx

	// newtxdata, si, ri := NewTxData(_accounts, *ethcl1, _ctx, r)
	// preTxBal := CalculateFunds(*ethcl1, _ctx, _accounts[si])
	// fmt.Printf("account %x has %v funds\n", _accounts[si].Address, preTxBal)
	// // unlock sender account
	// ks.Unlock(_accounts[si], passw0rd)

	// // fmt.Println("generating new transaction")

	// newtx := NewTx(newtxdata)
	// fmt.Println("signing transaction")
	// signedTx, _ := ks.SignTxWithPassphrase(_accounts[si], passw0rd, newtx, l1ChainID)
	// fmt.Println("signed")

	// fmt.Printf("L1 Tx: account %x with funds %v\n wants to send %v to account %x\n", _accounts[si].Address, preTxBal, newtxdata._amount, _accounts[ri].Address)

	// err = ethcl1.SendTransaction(_ctx, signedTx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// log.Fatal(err)
	// }
	
	if doDeposits {
		for len(zeros) != 0 {
			zeros = nil

			// Generate l1bridge.DispatchMessage inputs, deposit 1/10000 to l2
			disMsgData, si, _ := NewDmsgData(&ks, 100000, _accounts, *ethcl1, _ctx, r, l1ChainID)
			// fmt.Printf("Sender account: %v\n", _accounts[si].Address)
			// nnc, _ := ethcl1.NonceAt(_ctx, _accounts[si].Address, nil)
			fmt.Printf("account : %x || nonce: %v\n", _accounts[si].Address, disMsgData._nonce)
			//Send transaction with generated bind.TransactOpts
			tx, err := bridge.DispatchMessage(disMsgData._txOpts, disMsgData._to, disMsgData._fee, disMsgData._deadline, disMsgData._nonce, disMsgData._data)
			if err != nil {
				fmt.Println(err)
				// log.Fatal(err)
			}

			fmt.Printf("TxHash: %v\n", tx.Hash())
			TxHashes = append(TxHashes, tx.Hash())

			// fmt.Printf("TxCost: %v\n\n\n", tx.Cost())

			errr := ethcl1.SendTransaction(_ctx, tx)
			if errr != nil {
				fmt.Println(errr)
				// log.Fatal(err)
			}

			bal = GetBalances(_accounts, *ethcl1, *ethcl2, _ctx)
			for _, b := range bal {
				// fmt.Printf("account %x has %v funds in l2 and %v funds in l1\n", b.hexaddr, b.layer2Funds, b.layer1Funds)
				if b.layer2Funds.Cmp(big.NewInt(0)) == 0 {
					zeros = append(zeros, b.layer2Funds)
					// fmt.Println("true")
				}

			}
			// fmt.Println(zeros)
			// time.Sleep(3 * time.Second)
		}

		fmt.Printf("number of transactions sent: %v\n", len(TxHashes))
		fmt.Println(TxHashes)
	}



	if doTxs{
		ii := 0
		for true {
			newtxdata, si, ri := NewTxData(_accounts, *ethcl2, _ctx, r)
			newtx := NewTx(newtxdata)
			signedTx, err := ks.SignTxWithPassphrase(_accounts[si], passw0rd, newtx, l2ChainID)
			fmt.Printf("Chain ID: %v\n", l2ChainID)
			if err != nil {
				fmt.Println(err)
				// log.Fatal(err)
			}

			fmt.Printf("L2 Tx: account %x wants to send %v to account %x\n", _accounts[si].Address, newtxdata._amount, _accounts[ri].Address)
			fmt.Printf("Tx Hash: %v\n", signedTx.Hash())
			err = ethcl2.SendTransaction(_ctx, signedTx)
			if err != nil {
				fmt.Println(err)
				// log.Fatal(err)
			}

			time.Sleep(1 * time.Second)
			ii += 1
			fmt.Printf("Sent %v Txs so far\n", ii)
		}
	}
}
