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
	private := "cPQrREyy7VB5hjuZizRVjseRiHtsEt4q27BwwzHsdmN7hNqkQwXp"
	fromAddress := "mnansJbr3BFvfrW3TqhnygvTvYFQzHC4FF"
	toAddress := "msZJAsyxmmuxCLF58zSfa8R1XyHhQFG17Y"
	//
	api := cryptoapis.NewAPIClient("testnet", "2a454e24881ca117ca2201462c1e18691a15f9a5")

	rawTx, err := api.GetTransaction(context.Background(), fromAddress, toAddress, 1698000/1e8, 1000/1e8)
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

	//hash := &chainhash.Hash{}
	//hash.SetBytes([]byte("ee8d2f43b9169e50414fbf3a7150e90809dbfcd5c2ce011aed0e0b08f3b97361"))
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
	//caTx.TxIn[0].PreviousOutPoint = *(wire.NewOutPoint(hash, 0))

	fmt.Println()
	sourceAddress, err := btcutil.DecodeAddress(fromAddress, &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Println(err)
		return
	}

	sourcePkScript, _ := txscript.PayToAddrScript(sourceAddress)

	fmt.Println()

	//for _, txout := range caTx.TxOut {
	//	_, address, _, err := txscript.ExtractPkScriptAddrs(txout.PkScript, &chaincfg.TestNet3Params)
	//	if err != nil {
	//		fmt.Println("ERROR")
	//		fmt.Println(err)
	//		return
	//	}
	//}

	wif, err := btcutil.DecodeWIF(private)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	for index := range caTx.TxIn {
		sigScript, err := txscript.SignatureScript(caTx, index, sourcePkScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			fmt.Println(err)
			return
		}

		caTx.TxIn[index].SignatureScript = sigScript

		//flags := txscript.StandardVerifyFlags
		//vm, err := txscript.NewEngine(sourcePkScript, caTx, index, flags, nil, nil, 0)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//if err := vm.Execute(); err != nil {
		//	fmt.Println("HERE")
		//	fmt.Println(err)
		//	return
		//}
	}

	//var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	//sourceTx.Serialize(&unsignedTx)
	caTx.Serialize(&signedTx)

	fmt.Println(hex.EncodeToString(signedTx.Bytes()))
	signed := hex.EncodeToString(signedTx.Bytes())
	fmt.Println()
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
