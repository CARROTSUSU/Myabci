package main

import (
	"fmt"
	"log"

	"github.com/tendermint/tendermint/abci"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/libs/json"
)

// Struktur aplikasi RcpuCoin
type RcpuCoinApp struct {
	abci.BaseApplication
	balances map[string]int // Penyimpanan dalam ingatan untuk baki akaun
}

// Fungsi untuk mencipta aplikasi RcpuCoin baru
func NewRcpuCoinApp() *RcpuCoinApp {
	app := &RcpuCoinApp{
		balances: make(map[string]int),
	}

	// Menambah beberapa akaun permulaan dengan baki
	app.balances["address1"] = 1000
	app.balances["address2"] = 500

	return app
}

// Menerima dan memproses transaksi
func (app *RcpuCoinApp) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	var transfer struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}

	// Mendekod transaksi
	err := json.Unmarshal(tx, &transfer)
	if err != nil {
		return abci.ResponseDeliverTx{
			Code: 1,
			Log:  fmt.Sprintf("Error decoding transaction: %v", err),
		}
	}

	// Validasi jika penghantar ada baki yang mencukupi
	if app.balances[transfer.From] < transfer.Amount {
		return abci.ResponseDeliverTx{
			Code: 1,
			Log:  "Insufficient balance",
		}
	}

	// Mengemas kini baki akaun
	app.balances[transfer.From] -= transfer.Amount
	app.balances[transfer.To] += transfer.Amount

	// Log transaksi
	fmt.Printf("Transfer %d RCPU Coin from %s to %s\n", transfer.Amount, transfer.From, transfer.To)

	// Menyediakan respons untuk transaksi yang berjaya
	return abci.ResponseDeliverTx{
		Code: 0,
		Log:  "Transfer successful",
	}
}

// Fungsi untuk mendapatkan informasi aplikasi
func (app *RcpuCoinApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	// Menyediakan maklumat aplikasi
	return abci.ResponseInfo{
		Data:             "RCUCoin Blockchain",
		Validators:       []abci.Validator{},
		LatestBlockHash:  []byte("latest-block-hash"), // Gantikan dengan hash terkini jika perlu
		LatestAppHash:    []byte("latest-app-hash"),    // Gantikan dengan hash aplikasi terkini jika ada
	}
}

// Fungsi utama untuk menjalankan aplikasi
func main() {
	// Cipta aplikasi baru
	app := NewRcpuCoinApp()

	// Mulakan server ABCI untuk aplikasi ini pada port :26658
	server := abci.NewServer(":26658", app)
	log.Fatal(server.Start())
}
