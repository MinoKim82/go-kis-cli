package domestic

import (
	"fmt"

	"github.com/MinoKim82/go-kis-cli/pkg/auth"
	"github.com/MinoKim82/go-kis-cli/pkg/client"
)

type QuoteResponse struct {
	RtCd   string `json:"rt_cd"`
	MsgCd  string `json:"msg_cd"`
	Msg1   string `json:"msg1"`
	Output struct {
		StckPrpr   string `json:"stck_prpr"`    // 주식 현재가
		PrdyVrss   string `json:"prdy_vrss"`    // 전일 대비
		PrdyCtrt   string `json:"prdy_ctrt"`    // 전일 대비율
		AcmlVol    string `json:"acml_vol"`     // 누적 거래량
		AcmlTrPbmn string `json:"acml_tr_pbmn"` // 누적 거래 대금
		StckHgpr   string `json:"stck_hgpr"`    // 고가
		StckLwpr   string `json:"stck_lwpr"`    // 저가
		StckOpRC   string `json:"stck_oprc"`    // 시가
	} `json:"output"`
}

func GetQuote(stockCode string) (*QuoteResponse, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	token, err := auth.GetValidToken()
	if err != nil {
		return nil, err
	}

	var result QuoteResponse
	resp, err := c.Request("FHKST01010100", token).
		SetQueryParams(map[string]string{
			"FID_COND_MRKT_DIV_CODE": "J",
			"FID_INPUT_ISCD":         stockCode,
		}).
		SetResult(&result).
		Get("/uapi/domestic-stock/v1/quotations/inquire-price")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	if !client.IsSuccess(result.RtCd) {
		return nil, fmt.Errorf("API Error: [%s] %s", result.MsgCd, result.Msg1)
	}

	// Try to extract dynamic error message if output is empty and rt_cd is not 0 (but we checked above)
	return &result, nil
}
