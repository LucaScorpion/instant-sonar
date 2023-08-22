package sonar

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const apiContentType = "application/x-www-form-urlencoded"

type Client struct {
	Url      string
	username string
	password string
}

type tokenResponse struct {
	Token string `json:"token"`
}

func NewApiClient(url, username, password string) *Client {
	return &Client{
		Url:      url,
		username: username,
		password: password,
	}
}

func (client *Client) CreateProject(name, key string) {
	client.request(
		http.MethodPost,
		"/api/projects/create",
		createBody(map[string]string{
			"name":    name,
			"project": key,
		}),
	)
}

func (client *Client) CreateToken(projectKey string) string {
	res := client.request(
		http.MethodPost,
		"/api/user_tokens/generate",
		createBody(map[string]string{
			"name":       "Analyze " + projectKey,
			"projectKey": projectKey,
			"type":       "PROJECT_ANALYSIS_TOKEN",
		}),
	)

	var tokenRes tokenResponse
	resBytes, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(resBytes, &tokenRes)
	if err != nil {
		panic(err)
	}

	return tokenRes.Token
}

func (client *Client) DisableForceUserAuth() error {
	res := client.request(
		http.MethodPost,
		"/api/settings/set",
		createBody(map[string]string{
			"key":   "sonar.forceAuthentication",
			"value": "false",
		}),
	)

	if res.StatusCode != http.StatusNoContent {
		return errors.New("invalid sonarqube credentials")
	}
	return nil
}

func (client *Client) request(method, path string, body io.Reader) *http.Response {
	req, _ := http.NewRequest(method, client.Url+path, body)
	req.SetBasicAuth(client.username, client.password)
	req.Header.Set("Content-Type", apiContentType)
	res, _ := http.DefaultClient.Do(req)
	return res
}

func (client *Client) ProjectDashboardUrl(projectKey string) string {
	return client.Url + "/dashboard?id=" + projectKey
}

func createBody(body map[string]string) io.Reader {
	values := make(url.Values)
	for k, v := range body {
		values[k] = []string{v}
	}
	return strings.NewReader(values.Encode())
}
