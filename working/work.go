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
	private := "cQWYamWrmLb9iU8KVqZ9KizKs1MhmLM4FhH2Lw9YX9xM5iyLpHmt"
	fromAddress := "n2tRwYd6orVvU1yWxAeJBXvVyFbQKBT7gt"
	toAddress := "n2tRwYd6orVvU1yWxAeJBXvVyFbQKBT7gt"
	//
	api := cryptoapis.NewAPIClient("testnet", "2a454e24881ca117ca2201462c1e18691a15f9a5")

	rawTx, err := api.GetTransaction(context.Background(), fromAddress, toAddress, 0.003, 0.0001)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	fmt.Println(rawTx)

	//rawTx := "020000000142b9229bf67760ff6631553682bc9732a77e181fb164759abf19894d32369fc60000000000ffffffff02a0860100000000001976a91465e764fa399470b23d68138e66a3e216b156a33d88ac00350c00000000001976a914ea6a746899b49bb1c7d0229665e1d652129b942488ac00000000"
	//
	//fmt.Println("Raw tx: " + rawTx)
	////
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
	////
	//for _, txout := range caTx.TxOut {
	//	_, address, regSig, err := txscript.ExtractPkScriptAddrs(txout.PkScript, &chaincfg.TestNet3Params)
	//	if err != nil {
	//		fmt.Println("ERROR")
	//		fmt.Println(err)
	//	}
	//	fmt.Println(txout.Value)
	//	fmt.Println(address, regSig)
	//}

	fmt.Println()
	sourceTx := wire.NewMsgTx(wire.TxVersion)

	sourceTxIn := wire.NewTxIn(&caTx.TxIn[0].PreviousOutPoint, nil, nil)
	sourceTx.AddTxIn(sourceTxIn)
	sourceAddress, err := btcutil.DecodeAddress(fromAddress, &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Println(err)
		return
	}

	sourcePkScript, _ := txscript.PayToAddrScript(sourceAddress)
	//sourceTxHash := sourceTx.TxHash()
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	//prevOut := wire.NewOutPoint(&sourceTxHash, 0)
	for _, txIn := range caTx.TxIn {
		redeemTxIn := wire.NewTxIn(&txIn.PreviousOutPoint, nil, nil)
		redeemTx.AddTxIn(redeemTxIn)
	}

	fmt.Println()

	for _, txout := range caTx.TxOut {
		_, address, _, err := txscript.ExtractPkScriptAddrs(txout.PkScript, &chaincfg.TestNet3Params)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Println(err)
			return
		}

		destinationAddress, err := btcutil.DecodeAddress(address[0].String(), &chaincfg.TestNet3Params)
		if err != nil {
			fmt.Println(err)
			return
		}

		destinationPkScript, _ := txscript.PayToAddrScript(destinationAddress)
		redeemTxOut := wire.NewTxOut(txout.Value, destinationPkScript)
		redeemTx.AddTxOut(redeemTxOut)
	}

	wif, err := btcutil.DecodeWIF(private)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	for index := range caTx.TxIn {
		sigScript, err := txscript.SignatureScript(caTx, index, sourcePkScript, txscript.SigHashAll, wif.PrivKey, false)
		if err != nil {
			fmt.Println(err)
			return
		}

		caTx.TxIn[index].SignatureScript = sigScript

		flags := txscript.StandardVerifyFlags
		vm, err := txscript.NewEngine(sourcePkScript, caTx, index, flags, nil, nil, 0)
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

	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	caTx.Serialize(&signedTx)

	fmt.Println(hex.EncodeToString(signedTx.Bytes()))
	//signed := hex.EncodeToString(signedTx.Bytes())
	fmt.Println()
	//tx, err := api.SendSignedTransaction(context.Background(), signed)
	//if err != nil {
	//	fmt.Println("ERROR")
	//	fmt.Println(err)
	//	return
	//}
	//
	//
	//fmt.Println(tx)
	//sourceTxHash := sourceTx.TxHash()
	//redeemTx := wire.NewMsgTx(wire.TxVersion)
	//prevOut := wire.NewOutPoint(&sourceTxHash, 0)
	//redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	//redeemTx.AddTxIn(redeemTxIn)
	//redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	//redeemTx.AddTxOut(redeemTxOut)
	//sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	//if err != nil {
	//	return Transaction{}, err
	//}
	//redeemTx.TxIn[0].SignatureScript = sigScript
	//flags := txscript.StandardVerifyFlags
	//vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	//if err != nil {
	//	return Transaction{}, err
	//}
	//if err := vm.Execute(); err != nil {
	//	return Transaction{}, err
	//}
	//var unsignedTx bytes.Buffer
	//var signedTx bytes.Buffer
	//sourceTx.Serialize(&unsignedTx)
	//redeemTx.Serialize(&signedTx)

	//

	////
	////addresspubkey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	////fmt.Println(addresspubkey)
	//////
	////
	//sigScript, err := txscript.SignatureScript(sourceTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	//if err != nil {
	//	fmt.Println("ERROR")
	//	fmt.Println(err)
	//	return
	//}
	//
	//sourceTx.TxIn[0].SignatureScript = sigScript
	//flags := txscript.StandardVerifyFlags
	//vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, sourceTx, 0, flags, nil, nil, sourceTx.TxOut[0].Value)
	//if err != nil {
	//	fmt.Println("ERR")
	//	fmt.Println(err)
	//	return
	//}
	//if err := vm.Execute(); err != nil {
	//	fmt.Println("ERROR")
	//	fmt.Println(err)
	//	return
	//}
	//
	//var signedTx bytes.Buffer
	//
	//if err := sourceTx.Serialize(&signedTx); err != nil {
	//	fmt.Println("ERROR")
	//	fmt.Println(err)
	//	return
	//}
	//
	//signedRaw := hex.EncodeToString(signedTx.Bytes())
	//
	//fmt.Println()
	//fmt.Println(signedRaw)
	//
	//tx, err := api.SendSignedTransaction(context.Background(), signedRaw)
	//if err != nil {
	//	fmt.Println("ERROR")
	//	fmt.Println(err)
	//	return
	//}
	//
	//
	//fmt.Println(tx)

	//fmt.Println(wif)
	//fmt.Println(address)
	//fmt.Println(address.String())

}

func GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
}
