package sonar

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"instant-sonar/internal/docker"
)

const SonarqubeImage = "sonarqube"
const SonarScannerImage = "sonarsource/sonar-scanner-cli"
const SonarqubeOperationalMsg = "SonarQube is operational"

func CreateSonarQubeContainer(cli *docker.Client) string {
	res, err := cli.Cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: SonarqubeImage,
			ExposedPorts: nat.PortSet{
				"9000/tcp": {},
			},
		},
		nil,
		nil,
		nil,
		"sonarqube",
	)
	if err != nil {
		panic(err)
	}

	return res.ID
}

func CreateSonarScannerContainer(cli *docker.Client, hostUrl, projectKey, token, scanDir string, uid string) string {
	res, err := cli.Cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: SonarScannerImage,
			Env: []string{
				"SONAR_HOST_URL=" + hostUrl,
				"SONAR_SCANNER_OPTS=-Dsonar.projectKey=" + projectKey,
				"SONAR_TOKEN=" + token,
			},
			Volumes: map[string]struct{}{
				"/usr/src": {},
			},
			User: uid,
		},
		&container.HostConfig{
			AutoRemove: true,
			Binds: []string{
				scanDir + ":/usr/src",
			},
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
