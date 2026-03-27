package client

import (
	"fmt"
)

type HashResponse struct {
	Hash string `json:"HASH"`
}

// GenerateHashkey generates a security hash for POST requests required by KIS API
func (c *KISClient) GenerateHashkey(body map[string]interface{}) (string, error) {
	var result HashResponse

	resp, err := c.RestyClient.R().
		SetHeader("appkey", c.Profile.AppKey).
		SetHeader("appsecret", c.Profile.AppSecret).
		SetBody(body).
		SetResult(&result).
		Post("/uapi/hashkey")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("Hashkey generation failed HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	return result.Hash, nil
}
