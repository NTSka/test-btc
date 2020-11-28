package capi_eth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"test/eth/helpers"
)

const pth = `https://api.cryptoapis.io/v1/bc/eth`

type Client struct {
	net    string
	apiKey string
}

func NewAPIClient(net string, apiKey string) *Client {
	return &Client{
		net:    net,
		apiKey: apiKey,
	}
}

func (ac *Client) GetNonce(ctx context.Context, address string) (uint64, error) {
	target, err := joinPath(pth, ac.net, "/address/", address, "/nonce")
	if err != nil {
		return 0, err
	}

	result := struct {
		Payload struct {
			Nonce uint64 `json:"nonce"`
		} `json:"payload"`
	}{}

	if err := ac.doGETRequest(ctx, target, &result); err != nil {
		return 0, err
	}

	return result.Payload.Nonce, nil
}

func (ac *Client) GasPrice(ctx context.Context) (*big.Int, error) {
	target, err := joinPath(pth, ac.net, "/txs/fee")
	if err != nil {
		return nil, err
	}

	result := struct {
		Payload struct {
			Standard string `json:"standard"`
		} `json:"payload"`
	}{}

	if err := ac.doGETRequest(ctx, target, &result); err != nil {
		return nil, err
	}

	float, err := strconv.ParseFloat(result.Payload.Standard, 64)
	if err != nil {
		return nil, err
	}

	return helpers.Float64ToBigInt(float, 1e9), nil
}

func (ac *Client) Estimate(ctx context.Context, from string, to string, value float64, data string) (uint64, error) {
	target, err := joinPath(pth, ac.net, "/txs/gas")
	if err != nil {
		return 0, err
	}

	req := struct {
		FromAddress string  `json:"fromAddress"`
		ToAddress   string  `json:"toAddress"`
		Value       float64 `json:"value"`
		Data        string  `json:"data"`
	}{
		FromAddress: from,
		ToAddress:   to,
		Value:       value,
		Data:        data,
	}

	result := struct {
		Payload struct {
			GasLimit string `json:"gasLimit"`
		}
	}{}

	err = ac.doPOSTRequest(ctx, target, &req, &result)
	if err != nil {
		return 0, err
	}

	gasLimit, err := strconv.Atoi(result.Payload.GasLimit)
	if err != nil {
		return 0, err
	}

	return uint64(gasLimit), nil
}

func (ac *Client) BroadCastTx(ctx context.Context, hex string) (string, error) {
	target, err := joinPath(pth, ac.net, "/txs/push")
	if err != nil {
		return "", err
	}

	req := struct {
		Hex string `json:"hex"`
	}{
		Hex: hex,
	}

	result := struct {
		Payload struct {
			Txid string `json:"hex"`
		} `json:"payload"`
	}{}

	err = ac.doPOSTRequest(ctx, target, &req, &result)
	if err != nil {
		return "", err
	}
	return result.Payload.Txid, err
}

func (ac *Client) doPOSTRequest(ctx context.Context, target string, body interface{}, result interface{}) error {
	rawBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return ac.doRequestWithResponse(ctx, http.MethodPost, target, bytes.NewReader(rawBody), result)
}

func (ac *Client) doGETRequest(ctx context.Context, target string, result interface{}) error {
	return ac.doRequestWithResponse(ctx, http.MethodGet, target, nil, result)
}

func (ac *Client) doRequestWithResponse(ctx context.Context, method, target string, body io.Reader, output interface{}) error {
	res, err := ac.doRequest(ctx, method, target, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)

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

func (ac *Client) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
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
