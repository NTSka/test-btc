package main

import (
	"context"
	"fmt"
	"test/eth/eth_client"
	"test/eth/infura"
)

const private = ""
const to = ""
const depositData = "0xd0e30db0"
const infuraURL = ""
const caAPIKey = ""
const net = "ropsten"
const value = 0.001

func main() {
	infuraClient, err := infura.NewInfura(infuraURL)
	if err != nil {
		panic(err)
	}

	txHash, err := infuraClient.Send(context.Background(), private, to, value, depositData)
	if err != nil {
		panic(err)
	}

	fmt.Println("infura tx: ", txHash)

	caClient := eth_client.NewETHClient(net, caAPIKey)

	txHash, err = caClient.Send(context.Background(), private, to, value, depositData)
	if err != nil {
		panic(err)
	}

	fmt.Println("cryptoapis tx: ", txHash)
}
