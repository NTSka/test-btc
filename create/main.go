package main

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func main() {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}
	wif, err := btcutil.NewWIF(secret, &chaincfg.TestNet3Params, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	address, err := GetAddress(wif)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(wif.String())
	fmt.Println(address.EncodeAddress())
}

func GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
}
