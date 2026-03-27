package overseas

import (
	"fmt"

	"github.com/MinoKim82/go-kis-cli/config"
	"github.com/MinoKim82/go-kis-cli/pkg/auth"
	"github.com/MinoKim82/go-kis-cli/pkg/client"
)

type OrderResponse struct {
	RtCd   string `json:"rt_cd"`
	MsgCd  string `json:"msg_cd"`
	Msg1   string `json:"msg1"`
	Output struct {
		KRX_FWDG_ORD_ORGNO string `json:"KRX_FWDG_ORD_ORGNO"`
		ODNO               string `json:"ODNO"`
		ORD_TMD            string `json:"ORD_TMD"`
	} `json:"output"`
}

func executeOrder(isBuy bool, exchange string, symbol string, qty string, price string) (*OrderResponse, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	token, err := auth.GetValidToken()
	if err != nil {
		return nil, err
	}

	var trID string
	if config.EnvName == "mock" {
		if isBuy {
			trID = "VTTT1002U"
		} else {
			trID = "VTTT1006U"
		}
	} else {
		if isBuy {
			trID = "JTTT1002U"
		} else {
			trID = "JTTT1006U"
		}
	}

	orderDvsn := "00"
	if price == "0" || price == "" {
		// KIS requires specific handling for market orders in overseas, usually using a different code or leaving UNPR 0.
		// For simplicity, we assume "00" is limit order, and if 0, market order isn't fully supported without different conditions,
		// but we'll use "00" and price "0" and let the API reject it if invalid. Market order for overseas differs by exchange.
		price = "0"
	}

	body := map[string]interface{}{
		"CANO":            c.Profile.Cano,
		"ACNT_PRDT_CD":    c.Profile.PrdtCd,
		"OVRS_EXCG_CD":    exchange,
		"PDNO":            symbol,
		"ORD_QTY":         qty,
		"OVRS_ORD_UNPR":   price,
		"ORD_SVR_DVSN_CD": "0",
		"ORD_DVSN":        orderDvsn,
	}

	hash, err := c.GenerateHashkey(body)
	if err != nil {
		return nil, err
	}

	var result OrderResponse
	resp, err := c.Request(trID, token).
		SetHeader("hashkey", hash).
		SetBody(body).
		SetResult(&result).
		Post("/uapi/overseas-stock/v1/trading/order")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	if !client.IsSuccess(result.RtCd) {
		return nil, fmt.Errorf("API Error: [%s] %s", result.MsgCd, result.Msg1)
	}

	return &result, nil
}

func BuyOrder(exchange, symbol, qty, price string) (*OrderResponse, error) {
	return executeOrder(true, exchange, symbol, qty, price)
}

func SellOrder(exchange, symbol, qty, price string) (*OrderResponse, error) {
	return executeOrder(false, exchange, symbol, qty, price)
}
