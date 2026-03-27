package domestic

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
		Pdno        string `json:"pdno"`          // 종목번호
		PrdtName    string `json:"prdt_name"`     // 종목명
		HldgQty     string `json:"hldg_qty"`      // 보유수량
		Prpr        string `json:"prpr"`          // 현재가
		EvluPflsAmt string `json:"evlu_pfls_amt"` // 평가손익금액
		EvluPflsRt  string `json:"evlu_pfls_rt"`  // 평가손익율
	} `json:"output1"`
	Output2 []struct {
		TotEvluAmt string `json:"tot_evlu_amt"` // 총평가금액
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

	trID := "TTTC8434R"
	if config.EnvName == "mock" {
		trID = "VTTC8434R"
	}

	// KIS requires specific query parameters for balance inquiry
	params := map[string]string{
		"CANO":                  c.Profile.Cano,
		"ACNT_PRDT_CD":          c.Profile.PrdtCd,
		"AFHR_FLPR_YN":          "N",
		"OFL_YN":                "",
		"INQR_DVSN":             "01",
		"UNPR_DVSN":             "01",
		"FUND_STTL_ICLD_YN":     "N",
		"FNCG_AMT_AUTO_RDPT_YN": "N",
		"PRCS_DVSN":             "01",
		"CTX_AREA_FK100":        "",
		"CTX_AREA_NK100":        "",
	}

	var result BalanceResponse
	resp, err := c.Request(trID, token).
		SetQueryParams(params).
		SetResult(&result).
		Get("/uapi/domestic-stock/v1/trading/inquire-balance")

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
