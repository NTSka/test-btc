package eth_client

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"test/capi_eth"
	"test/eth/helpers"
	"time"
)

var chainIDS = map[string]*big.Int{
	"mainnet": big.NewInt(1),
	"ropsten": big.NewInt(3),
	"rinkeby": big.NewInt(4),
	"goerli":  big.NewInt(5),
}

type ETHClient struct {
	client  *capi_eth.Client
	chainId *big.Int
}

func NewETHClient(net string, apiKey string) *ETHClient {
	client := capi_eth.NewAPIClient(net, apiKey)

	chainId := chainIDS[net]

	return &ETHClient{client: client, chainId: chainId}
}

func (c *ETHClient) Send(ctx context.Context, private string, to string, value float64, data string) (string, error) {
	gasPrice, err := c.client.GasPrice(ctx)
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

	time.Sleep(time.Second)
	nonce, err := c.client.GetNonce(ctx, fromAddress.String())
	if err != nil {
		return "", err
	}

	fmt.Println("nonce: ", nonce)

	time.Sleep(time.Second)
	estimated, err := c.client.Estimate(ctx, fromAddress.String(), to, value, common.Bytes2Hex([]byte(data)))
	if err != nil {
		return "", err
	}

	fmt.Println("estimated CA: ", estimated)

	bigIntValue := helpers.Float64ToBigInt(value, 1e18)
	tx := types.NewTransaction(nonce, common.HexToAddress(to), bigIntValue, estimated, gasPrice, []byte(data))

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.chainId), privateKey)
	if err != nil {
		return "", err
	}

	ts := types.Transactions{signedTx}
	rawTx := hex.EncodeToString(ts.GetRlp(0))

	time.Sleep(time.Second)
	txHash, err := c.client.BroadCastTx(ctx, "0x"+rawTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
