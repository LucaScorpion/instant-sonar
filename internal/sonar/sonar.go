package sonar

import (
	"io"
	"net/http"
	"strings"
)

const SonarqubeImage = "sonarqube"
const SonarqubeOperationalMsg = "SonarQube is operational"

const apiContentType = "application/x-www-form-urlencoded"

type SonarQubeClient struct {
	url      string
	username string
	password string
}

func NewDefaultLocalSonarQubeClient() *SonarQubeClient {
	return &SonarQubeClient{
		url:      "http://localhost:9000",
		username: "admin",
		password: "admin",
	}
}

func (client *SonarQubeClient) CreateProject(name, key string) {
	client.request(
		http.MethodPost,
		"/api/projects/create",
		strings.NewReader("name="+name+"&project="+key),
	)
}

func (client *SonarQubeClient) request(method, path string, body io.Reader) *http.Response {
	req, _ := http.NewRequest(method, client.url+path, body)
	req.SetBasicAuth(client.username, client.password)
	req.Header.Set("Content-Type", apiContentType)
	res, _ := http.DefaultClient.Do(req)
	return res
}

func (client *SonarQubeClient) ProjectDashboardUrl(key string) string {
	return client.url + "/dashboard?id=" + key
}
