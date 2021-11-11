package github

import (
	"encoding/json"
	"fmt"
	"github-app-token/http/github/responses"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var client = &http.Client{}
var baseURL = os.Getenv("GITHUB_API_URL")

func formUrl(path string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(u).String(), err
}

func httpGet(path string, jwt string) ([]byte, error) {
	url, err := formUrl(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if jwt != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func httpPost(path string, jwt string) ([]byte, error) {
	url, err := formUrl(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	if jwt != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetApp(name string) (string, error) {
	raw, err := httpGet(fmt.Sprintf("apps/%s", name), "")
	if err != nil {
		return "", err
	}
	var app responses.App
	err = json.Unmarshal(raw, &app)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(app.Id), nil
}

func GetIntegrations(jwt string) (map[string]int, error) {
	out := map[string]int{}
	raw, err := httpGet("app/installations", jwt)
	if err != nil {
		return nil, err
	}
	var integrations responses.Integrations
	err = json.Unmarshal(raw, &integrations)
	if err != nil {
		return nil, err
	}
	for _, integration := range integrations {
		out[integration.Account.Login] = integration.Id
	}
	return out, nil
}

func GetAccessToken(jwt string, integration int) (string, error) {
	raw, err := httpPost(fmt.Sprintf("app/installations/%d/access_tokens", integration), jwt)
	if err != nil {
		return "", err
	}
	var accessToken responses.AccessToken
	err = json.Unmarshal(raw, &accessToken)
	if err != nil {
		return "", err
	}
	return accessToken.Token, nil
}
