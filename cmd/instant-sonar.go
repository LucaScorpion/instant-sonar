package main

import (
	"bufio"
	"fmt"
	"github.com/docker/go-connections/nat"
	"instant-sonar/internal"
	"instant-sonar/internal/docker"
	"instant-sonar/internal/sonar"
	"strings"
)

func main() {
	fmt.Println("Starting instant Sonar")
	cli := docker.NewDockerClient()
	defer cli.Close()

	contId := ""
	if cont, exists := docker.FindContainerByImageName(cli, sonar.SonarqubeImage); exists {
		contId = cont.ID
		fmt.Println("SonarQube container already exists (" + docker.ShortId(contId) + ")")

		if cont.State != docker.RunningState {
			fmt.Println("Starting SonarQube container")
			docker.StartContainer(cli, contId)
		}
	} else {
		fmt.Println("Pulling SonarQube image")
		docker.PullImage(cli, sonar.SonarqubeImage)
		fmt.Print("Creating SonarQube container")
		contId = docker.CreateContainer(cli, sonar.SonarqubeImage, "sonarqube", []nat.Port{"9000/tcp"})
		fmt.Println(" (" + docker.ShortId(contId) + ")")
		fmt.Println("Starting SonarQube container")
		docker.StartContainer(cli, contId)
	}

	fmt.Println("Waiting for SonarQube to be operational")
	out := docker.FollowContainerLogStream(cli, contId)
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

	sonarApi := sonar.NewDefaultLocalSonarQubeClient()

	fmt.Print("Creating project")
	projectKey := internal.RandomString(16)
	sonarApi.CreateProject(projectKey, projectKey)
	fmt.Println(" (" + projectKey + ")")

	fmt.Println("Project dashboard: " + sonarApi.ProjectDashboardUrl(projectKey))
}
