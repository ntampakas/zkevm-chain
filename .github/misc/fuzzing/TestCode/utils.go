package main

import (
	"context"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type ctx context.Context

type accs []accounts.Account

type netw0rks []netw0rk

type netw0rk struct {
	chainID *big.Int
	layer   string
	_url    string
}

type addr map[string]string

type txData struct {
	_to       common.Address
	_nonce    uint64
	_amount   *big.Int
	_gasLimit uint64
	_gasPrice *big.Int
	_data     []byte
}

type balance struct {
	hexaddr     common.Address
	layer1Funds *big.Int
	layer2Funds *big.Int
}

type balances []balance

type dMsgData struct {
	_txOpts   *bind.TransactOpts
	_to       common.Address
	_fee      *big.Int
	_deadline *big.Int
	_nonce    *big.Int
	_data     []byte
}

func NewDmsgData(ks *keystore.KeyStore, v int64, ac accs, ec ethclient.Client, c ctx, r *rand.Rand, chainid *big.Int) (dMsgData, int, int) {
	si := r.Intn(len(ac)) // sender's index (from accs slice)
	sender := ac[si]
	ri := r.Intn(len(ac)) // receiver's index (from accs slice)
	// n0nce, _ := ec.NonceAt(c, sender.Address, nil)
	// n0ncee := big.NewInt(int64(n0nce))
	receiver := ac[si]
	bal := CalculateFunds(ec, c, sender)
	txOpts, _ := bind.NewKeyStoreTransactorWithChainID(ks, sender, chainid)
	_nnc, _ := ec.NonceAt(c, sender.Address, nil)
	txOpts.Nonce = big.NewInt(int64(_nnc))
	// nnc := big.NewInt(int64(_nnc))
	// txOpts.Value = big.NewInt(v * params.GWei)
	txOpts.Value = new(big.Int).Div(bal, big.NewInt(v))
	txOpts.GasLimit = uint64(50000)
	// txOpts.GasTipCap = big.NewInt(255000 * params.GWei)
	txOpts.GasPrice, _ = ec.SuggestGasPrice(c)
	return dMsgData{
			txOpts,
			receiver.Address,
			big.NewInt(100 * params.GWei),
			big.NewInt(10000000000),
			txOpts.Nonce,
			[]byte{
				0x0,
			},
		},
		si,
		ri
}

func LoadAccounts(ksDir string) (accs, keystore.KeyStore) {
	ks := keystore.NewKeyStore(ksDir, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks.Accounts(), *ks
}

func CalculateFunds(ec ethclient.Client, c ctx, a accounts.Account) *big.Int {
	bal, _ := ec.BalanceAt(c, a.Address, nil)
	return bal
}

func NewTxData(ac accs, ec ethclient.Client, c ctx, r *rand.Rand) (txData, int, int) {
	si := r.Intn(len(ac)) // sender's index (from accs slice)
	sender := ac[si]
	ri := r.Intn(len(ac)) // receiver;s index (from accs slice)
	n0nce, _ := ec.NonceAt(c, sender.Address, nil)
	receiver := ac[ri]
	// _gasPrice, _ := ec.SuggestGasPrice(c)
	return txData{
			receiver.Address,
			n0nce,
			big.NewInt(1 * params.GWei),
			uint64(21000),
			big.NewInt(1 * params.GWei),
			// _gasPrice,
			[]byte{},
		},
		si,
		ri
}

func NewTx(newtxdata txData) *types.Transaction {
	tx := types.NewTransaction(
		newtxdata._nonce,
		newtxdata._to,
		newtxdata._amount,
		newtxdata._gasLimit,
		newtxdata._gasPrice,
		newtxdata._data,
	)
	return tx
}

func GetBalances(ac accs, ec1 ethclient.Client, ec2 ethclient.Client, c ctx) balances {
	_balances := balances{}
	for _, _acc := range ac {
		bal1 := CalculateFunds(ec1, c, _acc)
		// fmt.Println(_acc.Address, "layer1", bal1)
		bal2 := CalculateFunds(ec2, c, _acc)
		// fmt.Println(_acc.Address, "layer2", bal2)

		_balance := balance{
			_acc.Address,
			bal1,
			bal2,
		}
		// fmt.Println(_balance)
		_balances = append(_balances, _balance)
	}
	return _balances
}

// func KsSigner(ks keystore.KeyStore, a accounts.Account, passwd string, tx *types.Transaction, chainID *big.Int) (signerfn bind.SignerFn) {
// 	signedTx, _ := ks.SignTxWithPassphrase(a, passwd, tx, chainID)
// 	return signedTx
// }
