package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"instant-sonar/internal"
	"instant-sonar/internal/docker"
	"instant-sonar/internal/sonar"
	"strings"
)

func main() {
	fmt.Println("Starting instant Sonar")
	cli := docker.NewClient()
	defer cli.Close()

	qubeContId := ""
	if cont, exists := cli.FindContainerByImageName(sonar.SonarqubeImage); exists {
		qubeContId = cont.ID
		fmt.Println("SonarQube container already exists (" + docker.ShortId(qubeContId) + ")")

		if cont.State != docker.RunningState {
			fmt.Println("Starting SonarQube container")
			cli.StartContainer(qubeContId)
		}
	} else {
		fmt.Println("Pulling SonarQube image")
		cli.PullImage(sonar.SonarqubeImage)

		fmt.Print("Creating SonarQube container")
		qubeContId = sonar.CreateSonarQubeContainer(cli)
		fmt.Println(" (" + docker.ShortId(qubeContId) + ")")

		fmt.Println("Starting SonarQube container")
		cli.StartContainer(qubeContId)
	}

	fmt.Println("Waiting for SonarQube to be operational")
	out := cli.FollowContainerLogStream(qubeContId)
	bufOut := bufio.NewReader(out)

	for {
		bytes, _, err := bufOut.ReadLine()
		if err != nil {
			panic(err)
		}

		if strings.HasSuffix(string(bytes), sonar.SonarqubeOperationalMsg) {
			break
		}
	}
	out.Close()
	fmt.Println("SonarQube is operational")

	qubeWebUrl := "http://" + cli.GetContainerIp(qubeContId) + ":9000"
	sonarApi := sonar.NewApiClient(qubeWebUrl, "admin", "admin")
	sonarApi.DisableForceUserAuth()

	fmt.Print("Creating project")
	projectKey := internal.RandomString(16)
	sonarApi.CreateProject(projectKey, projectKey)
	fmt.Println(" (" + projectKey + ")")

	fmt.Print("Creating analysis token")
	token := sonarApi.CreateToken(projectKey)
	fmt.Println(" (" + token + ")")

	fmt.Println("Pulling Sonar Scanner image")
	cli.PullImage(sonar.SonarScannerImage)

	fmt.Print("Creating Sonar Scanner container")
	scanContId := sonar.CreateSonarScannerContainer(cli, sonarApi.Url, projectKey, token)
	fmt.Println(" (" + docker.ShortId(scanContId) + ")")

	fmt.Println("Starting analysis")
	cli.StartContainer(scanContId)
	cli.Cli.ContainerWait(context.Background(), scanContId, container.WaitConditionRemoved)
	fmt.Println("Done!")

	fmt.Println("Project dashboard: " + sonarApi.ProjectDashboardUrl(projectKey))
}
