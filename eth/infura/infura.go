package infura

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"test/eth/helpers"
)

type Infura struct {
	client *ethclient.Client
}

func NewInfura(url string) (*Infura, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Infura{client: client}, nil
}

func (i *Infura) Send(ctx context.Context, private string, to string, rawValue float64, data string) (string, error) {
	value := helpers.Float64ToBigInt(rawValue, 1e18)

	gasPrice, err := i.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	fmt.Println("gas price: ", gasPrice)

	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := i.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}

	fmt.Println("nonce: ", nonce)

	toAddress := common.HexToAddress(to)

	callMsg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Gas:      15000000,
		GasPrice: gasPrice,
		Value:    value,
		Data:     []byte(data),
	}

	estimated, err := i.client.EstimateGas(ctx, callMsg)
	if err != nil {
		return "", err
	}

	fmt.Println("estimated: ", estimated)

	tx := types.NewTransaction(nonce, toAddress, value, estimated, gasPrice, []byte(data))

	chainID, err := i.client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	if err := i.client.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}

	return signedTx.Hash().String(), nil
}
