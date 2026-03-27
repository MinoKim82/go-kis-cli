package domestic_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MinoKim82/go-kis-cli/config"
	"github.com/MinoKim82/go-kis-cli/pkg/domestic"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupConfig() {
	config.EnvName = "test"
	viper.Set("test.appkey", "mock-appkey")
	viper.Set("test.appsecret", "mock-appsecret")
}

func setupMockToken() {
	home, _ := os.UserHomeDir()
	cachePath := filepath.Join(home, ".kis-cli-token.json")
	mockToken := fmt.Sprintf(`{"access_token": "mock-token", "expires_at": "%s"}`, time.Now().Add(24*time.Hour).Format(time.RFC3339))
	os.WriteFile(cachePath, []byte(mockToken), 0600)
}

func TestGetQuote(t *testing.T) {
	setupConfig()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	setupMockToken()

	httpmock.RegisterResponder("GET", "https://openapivts.koreainvestment.com:29443/uapi/domestic-stock/v1/quotations/inquire-price",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "Bearer mock-token", req.Header.Get("authorization"))
			assert.Equal(t, "FHKST01010100", req.Header.Get("tr_id"))
			
			resp := domestic.QuoteResponse{
				RtCd:  "0",
				MsgCd: "MCA00000",
				Msg1:  "정상처리",
			}
			resp.Output.StckPrpr = "70000"
			resp.Output.PrdyVrss = "500"
			resp.Output.PrdyCtrt = "0.71"
			resp.Output.AcmlVol = "100000"

			jsonResp, _ := httpmock.NewJsonResponse(200, resp)
			return jsonResp, nil
		},
	)

	quote, err := domestic.GetQuote("005930")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	assert.NotNil(t, quote)
	assert.Equal(t, "0", quote.RtCd)
	assert.Equal(t, "70000", quote.Output.StckPrpr)
	assert.Equal(t, "100000", quote.Output.AcmlVol)
}
