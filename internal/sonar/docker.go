package sonar

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const SonarqubeImage = "sonarqube"
const SonarScannerImage = "sonarsource/sonar-scanner-cli"
const SonarqubeOperationalMsg = "SonarQube is operational"

func CreateSonarQubeContainer(cli *client.Client) string {
	res, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: SonarqubeImage,
			ExposedPorts: nat.PortSet{
				"9000/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"9000/tcp": []nat.PortBinding{{HostPort: "9000"}},
			},
		},
		nil,
		nil,
		"sonarqube",
	)
	if err != nil {
		panic(err)
	}

	return res.ID
}

func CreateSonarScannerContainer(cli *client.Client, hostUrl, projectKey, token string) string {
	// TODO: Volume mapping
	res, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: SonarScannerImage,
			Env: []string{
				"SONAR_HOST_URL=\"" + hostUrl + "\"",
				"SONAR_SCANNER_OPTS=\"-Dsonar.projectKey=" + projectKey + "\"",
				"SONAR_TOKEN=\"" + token + "\"",
			},
		},
		&container.HostConfig{
			AutoRemove:  true,
			NetworkMode: "host",
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		panic(err)
	}

	return res.ID
}
