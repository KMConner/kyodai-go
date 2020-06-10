package panda

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
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
	cookie []*http.Cookie
	url    *url.URL
	lt     string
	jar    *cookiejar.Jar
}

func getLoginForm() (*loginFormInfo, error) {
	req, err := http.NewRequest("GET", loginFormUrl, nil)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	cookie := resp.Cookies()

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
		cookie: cookie,
		url:    resp.Request.URL,
		lt:     string(matches[1]),
		jar:    jar,
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

	fmt.Println(formData.Encode())
	req, err := http.NewRequest("POST", formInfo.url.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range formInfo.cookie {
		req.AddCookie(cookie)
	}

	jar := formInfo.jar
	body, _ := httputil.DumpRequest(req, true)
	print(string(body))
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//return http.ErrUseLastResponse
			body, _ := httputil.DumpRequest(req, false)
			println(req.URL.String())
			println(string(body))
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ = httputil.DumpResponse(resp, true)
	fmt.Println(string(body))

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	println(bodyBytes)

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
