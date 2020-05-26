package kulasis

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type sessionLocation struct {
	JSession string
	Location string
}

type samlData struct {
	RelayState   string
	SamlResponse string
}

type cookieLocation struct {
	Cookie   string
	Location string
}

func getSessionId() (id string, loc string, err error) {
	const sessionCompleteUrl = "https://www.k.kyoto-u.ac.jp/api/app/v1/auth/get_j_session_complete"
	resp, e := http.Get(sessionCompleteUrl)
	if e != nil {
		return "", "", e
	}

	defer func() { _ = resp.Body.Close() }()

	bodyBytes, e := ioutil.ReadAll(resp.Body)

	if e != nil {
		return "", "", e
	}
	var sessionLoc sessionLocation
	e = json.Unmarshal(bodyBytes, &sessionLoc)
	if e != nil {
		return "", "", e
	}
	return sessionLoc.JSession, sessionLoc.Location, nil
}

func getLogInPage(cookie []*http.Cookie) error {
	const loginPageGetUrl = "https://www.k.kyoto-u.ac.jp/secure/student/shibboleth_account_list?keep=true"
	req, e := http.NewRequest(http.MethodGet, loginPageGetUrl, nil)
	if e != nil {
		return e
	}
	for _, c := range cookie {
		req.AddCookie(c)
	}

	client := http.DefaultClient
	resp, e := client.Do(req)

	if e != nil {
		return e
	}

	defer func() { _ = resp.Body.Close() }()
	return nil
}

func postLogin(url string, id string, pass string, sessionId string) (*samlData, error) {
	cookie := &http.Cookie{
		Domain: "authidp1.iimc.kyoto-u.ac.jp",
		Name:   "JSESSIONID",
		Value:  sessionId}
	dataStr := "j_username=" + id + "&j_password=" + pass + "&_eventId_proceed="

	req, e := http.NewRequest("POST", url, bytes.NewBufferString(dataStr))

	if e != nil {
		return nil, e
	}
	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.DefaultClient
	resp, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := extractSamlData(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func arrayToMap(tokens *html.Tokenizer) map[string]string {
	ret := make(map[string]string)
	for {
		key, value, hasNext := tokens.TagAttr()
		ret[string(key)] = string(value)
		if !hasNext {
			break
		}
	}
	return ret
}

func extractSamlData(reader io.Reader) (*samlData, error) {
	var state *string = nil
	var saml *string = nil

	tokenizer := html.NewTokenizer(reader)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}

		if tt != html.SelfClosingTagToken {
			continue
		}

		name, _ := tokenizer.TagName()
		if strings.ToLower(string(name)) != "input" {
			continue
		}
		attrs := arrayToMap(tokenizer)
		inputName, nameOk := attrs["name"]
		if !nameOk {
			continue
		}
		inputValue, valueOk := attrs["value"]
		if !valueOk {
			continue
		}
		switch inputName {
		case "RelayState":
			state = &inputValue
		case "SAMLResponse":
			saml = &inputValue
		default:
			continue
		}
	}

	if state == nil || saml == nil {
		return nil, errors.New("ERROR OCCURRED WHILE SIGN IN")
	}

	return &samlData{
		RelayState:   *state,
		SamlResponse: *saml,
	}, nil
}

func postSaml(saml *samlData) (cookie *http.Cookie, location string, err error) {
	const samlPostUrl = "https://www.k.kyoto-u.ac.jp/api/app/v1/auth/get_shibboleth_session"
	data := url.Values{}
	data.Set("RelayState", saml.RelayState)
	data.Set("SAMLResponse", saml.SamlResponse)
	data.Set("requestUrl", "https://www.k.kyoto-u.ac.jp/Shibboleth.sso/SAML2/POST")

	req, e := http.NewRequest("POST", samlPostUrl, strings.NewReader(data.Encode()))

	if e != nil {
		return nil, "", e
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return nil, "", e
	}

	bodyBytes, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, "", e
	}

	var cookieLoc cookieLocation
	e = json.Unmarshal(bodyBytes, &cookieLoc)
	if e != nil {
		return nil, "", e
	}

	keyValue := strings.Split(cookieLoc.Cookie, "=")
	if len(keyValue) < 2 {
		return nil, "", errors.New("ERROR WHEN SIGN IN")
	}
	retCookie := &http.Cookie{
		Name:   keyValue[0],
		Value:  keyValue[1],
		Domain: "www.k.kyoto-u.ac.jp",
	}

	return retCookie, cookieLoc.Location, nil
}

func getToken(cookie *http.Cookie) (*Info, error) {
	req, e := http.NewRequest("GET", "https://www.k.kyoto-u.ac.jp/secure/student/shibboleth_account_list?keep=true", nil)
	if e != nil {
		return nil, e
	}
	req.AddCookie(cookie)
	req.AddCookie(&http.Cookie{Name: "cserver", Value: "ku_europa", Domain: "www.k.kyoto-u.ac.jp"})

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return nil, e
	}

	bodyBytes, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	var info Info
	e = json.Unmarshal(bodyBytes, &info)
	if e != nil {
		return nil, e
	}

	return &info, nil
}

func SignIn(userName string, password string) (*Info, error) {
	sessionId, location, err := getSessionId()
	if err != nil {
		return nil, err
	}
	sessionId = strings.Split(sessionId, "=")[1]
	cookies := []*http.Cookie{{
		Domain: "authidp1.iimc.kyoto-u.ac.jp",
		Name:   "JSESSIONID",
		Value:  sessionId}, {
		Domain: "authidp1.iimc.kyoto-u.ac.jp",
		Name:   "cserver",
		Value:  "ku_europa"}}

	err = getLogInPage(cookies)
	if err != nil {
		return nil, err
	}

	samlData, err := postLogin(location, userName, password, sessionId)
	if err != nil {
		return nil, err
	}

	c, _, err := postSaml(samlData)
	if err != nil {
		return nil, err
	}

	info, err := getToken(c)
	if err != nil {
		return nil, err
	}

	return info, nil
}
