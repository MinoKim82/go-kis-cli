package overseas

import (
	"fmt"

	"github.com/MinoKim82/go-kis-cli/config"
	"github.com/MinoKim82/go-kis-cli/pkg/auth"
	"github.com/MinoKim82/go-kis-cli/pkg/client"
)

type BalanceResponse struct {
	RtCd    string `json:"rt_cd"`
	MsgCd   string `json:"msg_cd"`
	Msg1    string `json:"msg1"`
	Output1 []struct {
		OvrsPdno     string `json:"ovrs_pdno"`      // 해외종목번호
		OvrsItemName string `json:"ovrs_item_name"` // 해외종목명
		OvrsCblcQty  string `json:"ovrs_cblc_qty"`  // 해외잔고수량
		NowPric2     string `json:"now_pric2"`      // 현재가격
		EvluPflsAmt  string `json:"evlu_pfls_amt"`  // 평가손익금액
		EvluPflsRt   string `json:"evlu_pfls_rt"`   // 평가손익율
	} `json:"output1"`
	Output2 struct {
		FrcrEvluTamt string `json:"frcr_evlu_tamt"` // 외화평가총금액
	} `json:"output2"`
}

func GetBalance() (*BalanceResponse, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	token, err := auth.GetValidToken()
	if err != nil {
		return nil, err
	}

	trID := "JTTT3012R"
	if config.EnvName == "mock" {
		trID = "VTTS3012R"
	}

	params := map[string]string{
		"CANO":           c.Profile.Cano,
		"ACNT_PRDT_CD":   c.Profile.PrdtCd,
		"OVRS_EXCG_CD":   "NASD", // default nasdaq or all
		"TR_CRCY_CD":     "USD",
		"CTX_AREA_FK200": "",
		"CTX_AREA_NK200": "",
	}

	var result BalanceResponse
	resp, err := c.Request(trID, token).
		SetQueryParams(params).
		SetResult(&result).
		Get("/uapi/overseas-stock/v1/trading/inquire-balance")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	// Wait, KIS uses different RtCd handling sometimes but '0' is universal success.
	if !client.IsSuccess(result.RtCd) {
		return nil, fmt.Errorf("API Error: [%s] %s", result.MsgCd, result.Msg1)
	}

	return &result, nil
}
