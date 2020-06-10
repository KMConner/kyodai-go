package panda

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

const loginFormUrl = "https://panda.ecs.kyoto-u.ac.jp/portal/login"

const portalUrl = "https://panda.ecs.kyoto-u.ac.jp/portal"

type AuthInfo struct {
	jar    *cookiejar.Jar
	client *http.Client
}

type loginFormInfo struct {
	url *url.URL
	lt  string
	jar *cookiejar.Jar
}

func getLoginForm() (*loginFormInfo, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{Jar: jar}

	resp, err := client.Get(loginFormUrl)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(`<input type="hidden" name="lt" value="(.*?)".*?>`)
	if err != nil {
		return nil, err
	}

	matches := regex.FindSubmatch(bodyBytes)
	if len(matches) != 2 {
		return nil, errors.New("FAILED_TO_PARSE_LOGIN_FORM")
	}

	return &loginFormInfo{
		url: resp.Request.URL,
		lt:  string(matches[1]),
		jar: jar,
	}, nil
}

func postLogin(id string, pass string, formInfo *loginFormInfo) (*AuthInfo, error) {
	formData := url.Values{}
	formData.Set("_eventId", "submit")
	formData.Set("execution", "e1s1")
	formData.Set("lt", formInfo.lt)
	formData.Set("password", pass)
	formData.Set("username", id)
	formData.Set("submit", "ログイン")

	jar := formInfo.jar
	client := &http.Client{
		Jar: jar,
	}
	resp, err := client.Post(formInfo.url.String(),
		"application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))

	if err != nil {
		return nil, err
	}

	if resp.Request.URL.String() != portalUrl {
		return nil, errors.New("INVALID_ID_OR_PASSWORD")
	}

	return &AuthInfo{
		jar: jar,
	}, nil
}

func SignIn(id string, pass string) (*AuthInfo, error) {
	form, err := getLoginForm()
	if err != nil {
		return nil, err
	}
	return postLogin(id, pass, form)
}
