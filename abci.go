package main

import (
	"fmt"
	"log"

	"github.com/tendermint/tendermint/abci"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/libs/json"
)

type RcpuCoinApp struct {
	abci.BaseApplication
}

func NewRcpuCoinApp() *RcpuCoinApp {
	return &RcpuCoinApp{}
}

// Menerima transaksi
func (app *RcpuCoinApp) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	var transfer struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}

	err := json.Unmarshal(tx, &transfer)
	if err != nil {
		return abci.ResponseDeliverTx{
			Code: 1,
			Log:  fmt.Sprintf("Error decoding transaction: %v", err),
		}
	}

	// Logik untuk memindahkan RCPU Coin (contoh: update akaun penerima)
	fmt.Printf("Transfer %d RCPU Coin from %s to %s\n", transfer.Amount, transfer.From, transfer.To)

	// Berikan response
	return abci.ResponseDeliverTx{
		Code: 0,
		Log:  "Transfer successful",
	}
}

// Fungsi untuk aplikasi blok
func (app *RcpuCoinApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{
		Data:             "RCUCoin Blockchain",
		Validators:       []abci.Validator{},
		LatestBlockHash:  []byte{},
		LatestAppHash:    []byte{},
	}
}

func main() {
	app := NewRcpuCoinApp()

	// Jalankan server ABCI
	server := abci.NewServer(":26658", app)
	log.Fatal(server.Start())
}
