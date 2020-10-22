package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"test/cryptoapis"
)

func main() {
	private := "cQDYDrYjmJfXozzUohJdwk9nBHvChBJeLeC2RbdVvhioqqnX6ot5"
	fromAddress := "mnansJbr3BFvfrW3TqhnygvTvYFQzHC4FF"
	toAddress := "mnansJbr3BFvfrW3TqhnygvTvYFQzHC4FF"

	api := cryptoapis.NewAPIClient("testnet", "2a454e24881ca117ca2201462c1e18691a15f9a5")

	amount := 0.001
	fee := 0.0001

	rawTx, err := api.GetTransaction(context.Background(), fromAddress, toAddress, amount, fee)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	res, err := hex.DecodeString(rawTx)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	caTx := wire.NewMsgTx(wire.TxVersion)

	if err := caTx.Deserialize(bytes.NewReader(res)); err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	sourceAddress, err := btcutil.DecodeAddress(fromAddress, &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Println(err)
		return
	}

	sourcePkScript, _ := txscript.PayToAddrScript(sourceAddress)

	wif, err := btcutil.DecodeWIF(private)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	for _, txIn := range caTx.TxIn {
		fmt.Println(txIn.PreviousOutPoint.String())
		redeemTxIn := wire.NewTxIn(&txIn.PreviousOutPoint, nil, nil)
		redeemTx.AddTxIn(redeemTxIn)
	}

	destinationAddress, err := btcutil.DecodeAddress(toAddress, &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Println(err)
		return
	}
	destinationPkScript, _ := txscript.PayToAddrScript(destinationAddress)

	total := int64(0)
	for _, txOut := range caTx.TxOut {
		total += txOut.Value
	}

	outputsCount := float64(3)
	eachOutputAmount := (amount / outputsCount) * 1e8
	for i := 0; i < int(outputsCount); i++ {
		redeemTx.AddTxOut(wire.NewTxOut(int64(eachOutputAmount), destinationPkScript))
	}

	redeemTx.AddTxOut(wire.NewTxOut(total-int64(fee*1e8+eachOutputAmount*outputsCount), destinationPkScript))

	for index := range redeemTx.TxIn {
		sigScript, err := txscript.SignatureScript(redeemTx, index, sourcePkScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			fmt.Println(err)
			return
		}

		redeemTx.TxIn[index].SignatureScript = sigScript

		flags := txscript.StandardVerifyFlags
		vm, err := txscript.NewEngine(sourcePkScript, redeemTx, index, flags, nil, nil, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := vm.Execute(); err != nil {
			fmt.Println("HERE")
			fmt.Println(err)
			return
		}
	}

	//var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	//source
	//Tx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)

	signed := hex.EncodeToString(signedTx.Bytes())
	fmt.Println(signed)
	tx, err := api.SendSignedTransaction(context.Background(), signed)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	fmt.Println(tx)
}

func GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
}
