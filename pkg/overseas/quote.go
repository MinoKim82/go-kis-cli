package overseas

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
		Sign   string `json:"sign"`   // 대비기호
		Diff   string `json:"diff"`   // 전일대비
		Rate   string `json:"rate"`   // 등락율
		Last   string `json:"last"`   // 현재가
		Pclose string `json:"pclose"` // 전일종가
		Tvol   string `json:"tvol"`   // 거래량
		Tamt   string `json:"tamt"`   // 거래대금
	} `json:"output"`
}

func GetQuote(exchange string, symbol string) (*QuoteResponse, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	token, err := auth.GetValidToken()
	if err != nil {
		return nil, err
	}

	var result QuoteResponse
	resp, err := c.Request("HHDFS76200200", token).
		SetQueryParams(map[string]string{
			"AUTH": "",
			"EXCD": exchange, // NAS, NYS, AMS, etc
			"SYMB": symbol,
		}).
		SetResult(&result).
		Get("/uapi/overseas-price/v1/quotations/price")

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
