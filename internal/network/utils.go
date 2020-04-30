package network

import (
	"encoding/json"
	"github.com/KMConner/kyodai-go/internal/auth"
	"io/ioutil"
	"net/http"
	"net/url"
)

func AccessWithToken(accessUrl url.URL, credential *auth.Info, data interface{}) error {
	query := accessUrl.Query()
	query.Add("accessToken", credential.AccessToken)
	query.Add("account", credential.Account)
	accessUrl.RawQuery = query.Encode()

	resp, err := http.Get(accessUrl.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	return err
}
