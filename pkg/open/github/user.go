package github

import (
	"encoding/json"
	"errors"
	"net/http"
)

// https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/#3-use-the-access-token-to-access-the-api
func (c *Client) GetUserInfo(accessToken string) (resp UserInfo, err error) {
	response, err := c.RestyClient.R().SetHeader("Authorization", "token "+accessToken).Get(ApiUser)
	if err != nil {
		return
	}

	if response.StatusCode() != http.StatusOK {
		err = errors.New("status code is not 200")
		return
	}

	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return
	}
	return
}
