package cryptoapis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const pth = `https://api.cryptoapis.io/v1/bc/btc`

type APIClient struct {
	net    string
	apiKey string
}

func NewAPIClient(net string, apiKey string) *APIClient {
	return &APIClient{
		net:    net,
		apiKey: apiKey,
	}
}

func (ac *APIClient) SendSignedTransaction(ctx context.Context, hex string) (string, error) {
	target, err := joinPath(pth, ac.net, "/txs/send")
	if err != nil {
		return "", err
	}

	req := SignedTransaction{Hex: hex}

	result := struct {
		Payload struct {
			Txid string `json:"txid"`
		} `json:"payload"`
	}{}

	err = ac.doPOSTRequest(ctx, target, &req, &result)
	if err != nil {
		return "", err
	}
	return result.Payload.Txid, err
}

func (ac *APIClient) GetTransaction(ctx context.Context, from, to string, amount float64, fee float64) (string, error) {
	target, err := joinPath(pth, ac.net, "/txs/create")
	if err != nil {
		return "", err
	}

	req := Tx{
		Inputs: []Destination{{
			Address: from,
			Value:   amount,
		}},
		Outputs: []Destination{{
			Address: to,
			Value:   amount,
		}},
		Fee: Destination{
			Address: from,
			Value:   fee,
		},
	}

	result := struct {
		Payload struct {
			Hex string `json:"hex"`
		} `json:"payload"`
	}{}
	err = ac.doPOSTRequest(ctx, target, &req, &result)
	if err != nil {
		return "", err
	}

	return result.Payload.Hex, err
}

func (ac *APIClient) doGETRequest(ctx context.Context, target string, result interface{}) error {
	return ac.doRequestWithResponse(ctx, http.MethodGet, target, nil, result)
}

func (ac *APIClient) doPOSTRequest(ctx context.Context, target string, body interface{}, result interface{}) error {
	rawBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return ac.doRequestWithResponse(ctx, http.MethodPost, target, bytes.NewReader(rawBody), result)
}

func (ac *APIClient) doRequestWithResponse(ctx context.Context, method, target string, body io.Reader, output interface{}) error {
	res, err := ac.doRequest(ctx, method, target, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(rawBody))
	switch {
	case err != nil:
		return err
	case res.StatusCode != http.StatusOK:
		return err
	case output == nil:
		return nil
	}

	if t, ok := output.(*string); ok {
		*t = string(rawBody)
		return nil
	}

	err = json.Unmarshal(rawBody, &output)
	return err
}

func (ac *APIClient) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-API-Key", ac.apiKey)
	req.Header.Add("Content-Type", "application/json")
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func joinPath(basePath string, paths ...string) (string, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(append([]string{u.Path}, paths...)...)
	return u.String(), nil
}

type FeeResult struct {
	Min               QuotedFloat64 `json:"min"`
	Max               QuotedFloat64 `json:"max"`
	Average           QuotedFloat64 `json:"average"`
	MinFeePerByte     QuotedFloat64 `json:"min_fee_per_byte"`
	AverageFeePerByte QuotedFloat64 `json:"average_fee_per_byte"`
	MaxFeePerByte     QuotedFloat64 `json:"max_fee_per_byte"`
	Unit              string        `json:"unit"`
}

type QuotedFloat64 float64

func (q *QuotedFloat64) UnmarshalJSON(b []byte) error {
	trimmed := bytes.Trim(b, "\"")
	res, err := strconv.ParseFloat(string(trimmed), 64)
	*q = QuotedFloat64(res)
	return err
}

type SignedTransaction struct {
	Hex string `json:"hex"`
}

type SendMoney struct {
	CreateTx Tx       `json:"createTx"`
	Wifs     []string `json:"wifs"`
}

type Destination struct {
	Address string  `json:"address"`
	Value   float64 `json:"value"`
}

type Tx struct {
	Inputs  []Destination `json:"inputs"`
	Outputs []Destination `json:"outputs"`
	Fee     Destination   `json:"fee"`
}

type BlockData struct {
	Hash              string   `json:"hash"`
	Strippedsize      int      `json:"strippedsize"`
	Size              int      `json:"size"`
	Weight            int      `json:"weight"`
	Height            int      `json:"height"`
	Version           int      `json:"version"`
	VersionHex        string   `json:"versionHex"`
	Merkleroot        string   `json:"merkleroot"`
	Datetime          string   `json:"datetime"`
	Mediantime        string   `json:"mediantime"`
	Nonce             int      `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	Previousblockhash string   `json:"previousblockhash"`
	Nextblockhash     string   `json:"nextblockhash"`
	Transactions      int      `json:"transactions"`
	Tx                []string `json:"tx"`
	Confirmations     int      `json:"confirmations"`
	Timestamp         int      `json:"timestamp"`
}

type TransactionData struct {
	Txid          string `json:"txid"`
	Hash          string `json:"hash"`
	Index         int    `json:"index"`
	Version       int    `json:"version"`
	Size          int    `json:"size"`
	Vsize         int    `json:"vsize"`
	Locktime      int    `json:"locktime"`
	Time          string `json:"time"`
	Blockhash     string `json:"blockhash"`
	Blockheight   int    `json:"blockheight"`
	Blocktime     string `json:"blocktime"`
	Timestamp     int    `json:"timestamp"`
	Confirmations int    `json:"confirmations"`
	Txins         []struct {
		Txout     string   `json:"txout"`
		Vout      int      `json:"vout"`
		Amount    string   `json:"amount"`
		Addresses []string `json:"addresses"`
		Script    struct {
			Asm string `json:"asm"`
			Hex string `json:"hex"`
		} `json:"script"`
		Votype string `json:"votype"`
	} `json:"txins"`
	Txouts []struct {
		Amount    string   `json:"amount"`
		Type      string   `json:"type"`
		Spent     bool     `json:"spent"`
		Addresses []string `json:"addresses"`
		Script    struct {
			Asm     string `json:"asm"`
			Hex     string `json:"hex"`
			Reqsigs int    `json:"reqsigs"`
		} `json:"script"`
	} `json:"txouts"`
}
