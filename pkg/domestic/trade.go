package domestic

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
		ODNO               string `json:"ODNO"`    // 주문 번호
		ORD_TMD            string `json:"ORD_TMD"` // 주문 시간
	} `json:"output"`
}

// executeOrder is a helper function to fire Buy or Sell orders
func executeOrder(isBuy bool, stockCode string, qty string, price string) (*OrderResponse, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	token, err := auth.GetValidToken()
	if err != nil {
		return nil, err
	}

	// Determine TR_ID
	var trID string
	if config.EnvName == "mock" {
		if isBuy {
			trID = "VTTC0802U"
		} else {
			trID = "VTTC0801U"
		}
	} else {
		if isBuy {
			trID = "TTTC0802U"
		} else {
			trID = "TTTC0801U"
		}
	}

	// 00: 지정가, 01: 시장가. If price is 0, we assume market price
	orderDvsn := "00"
	if price == "0" || price == "" {
		orderDvsn = "01" // Market
		price = "0"
	}

	body := map[string]interface{}{
		"CANO":         c.Profile.Cano,
		"ACNT_PRDT_CD": c.Profile.PrdtCd,
		"PDNO":         stockCode,
		"ORD_DVSN":     orderDvsn,
		"ORD_QTY":      qty,
		"ORD_UNPR":     price,
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
		Post("/uapi/domestic-stock/v1/trading/order-cash")

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

func BuyOrder(stockCode, qty, price string) (*OrderResponse, error) {
	return executeOrder(true, stockCode, qty, price)
}

func SellOrder(stockCode, qty, price string) (*OrderResponse, error) {
	return executeOrder(false, stockCode, qty, price)
}
