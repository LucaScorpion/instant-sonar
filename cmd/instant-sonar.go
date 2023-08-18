package main

import (
	"fmt"
	"instant-sonar/internal/docker"
	"instant-sonar/internal/sonar"
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
		contId = docker.CreateContainer(cli, sonar.SonarqubeImage, "sonarqube")
		fmt.Println(" (" + docker.ShortId(contId) + ")")
		fmt.Println("Starting SonarQube container")
		docker.StartContainer(cli, contId)
	}
}
