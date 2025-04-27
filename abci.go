package main

import (
	"fmt"
	"github.com/tendermint/tendermint/abci"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/libs/log"
)

type RCUCOINApplication struct {
	// Simpan data yang diperlukan untuk transaksi dan token
	balances map[string]int64
}

// Fungsi Info untuk mendapatkan metadata aplikasi
func (app *RCUCOINApplication) Info(req abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{
		Data: "RCUCOIN ABCI App",
	}
}

// Fungsi DeliverTx untuk memproses transaksi
func (app *RCUCOINApplication) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	var sender, receiver string
	var amount int64
	// Parse transaksi dalam bentuk format yang disarankan (contoh: sender, receiver, amount)
	fmt.Sscanf(string(tx), "%s %s %d", &sender, &receiver, &amount)

	// Semak apakah saldo penghantar mencukupi
	if app.balances[sender] < amount {
		return abci.ResponseDeliverTx{
			Code: abci.CodeTypeUnauthorized,
			Log:  "Insufficient funds",
		}
	}

	// Pemindahan jumlah dari penghantar ke penerima
	app.balances[sender] -= amount
	app.balances[receiver] += amount

	return abci.ResponseDeliverTx{
		Code: abci.CodeTypeOK,
	}
}

// Fungsi CheckTx untuk memeriksa transaksi sebelum dihantar
func (app *RCUCOINApplication) CheckTx(tx []byte) abci.ResponseCheckTx {
	var sender, receiver string
	var amount int64
	fmt.Sscanf(string(tx), "%s %s %d", &sender, &receiver, &amount)

	// Periksa sama ada jumlah adalah positif
	if amount <= 0 {
		return abci.ResponseCheckTx{
			Code: abci.CodeTypeInvalidRequest,
			Log:  "Amount must be positive",
		}
	}

	return abci.ResponseCheckTx{
		Code: abci.CodeTypeOK,
	}
}

// Fungsi Commit untuk komit perubahan ke dalam blockchain
func (app *RCUCOINApplication) Commit() abci.ResponseCommit {
	// Simpan perubahan yang dibuat (contoh: menggunakan hash)
	return abci.ResponseCommit{}
}

// Fungsi InitChain untuk inisialisasi blockchain
func (app *RCUCOINApplication) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
	// Inisialisasi blockchain dan tetapkan salinan pertama bagi token
	app.balances = make(map[string]int64)
	app.balances["genesis"] = 1000000 // Genesis block mempunyai 1 juta RCUCOIN
	return abci.ResponseInitChain{}
}

// Fungsi BeginBlock untuk permulaan blok
func (app *RCUCOINApplication) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// Fungsi EndBlock untuk penutupan blok
func (app *RCUCOINApplication) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

// Fungsi Query untuk query status
func (app *RCUCOINApplication) Query(req abci.RequestQuery) abci.ResponseQuery {
	// Query untuk melihat saldo
	if req.Path == "/balance" {
		address := string(req.Data)
		balance := app.balances[address]
		return abci.ResponseQuery{
			Code: 0,
			Log:  fmt.Sprintf("Balance of %s is %d", address, balance),
		}
	}

	return abci.ResponseQuery{
		Code: abci.CodeTypeUnknownRequest,
		Log:  "Unknown query path",
	}
}

func main() {
	// Memulakan aplikasi ABCI dengan server soket
	app := &RCUCOINApplication{}
	server := abci.NewSocketServer(app, "tcp://127.0.0.1:26658")

	// Mulakan server
	if err := server.Start(); err != nil {
		fmt.Println("Error starting ABCI server:", err)
	}
}
