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
		fmt.Println("SonarQube container already exists")
		contId = cont.ID

		if cont.State != "running" {
			fmt.Println("Starting SonarQube container")
			docker.StartContainer(cli, contId)
		}
	} else {
		fmt.Println("Pulling SonarQube image")
		docker.PullImage(cli, sonar.SonarqubeImage)
		fmt.Println("Creating SonarQube container")
		contId = docker.CreateContainer(cli, sonar.SonarqubeImage, "sonarqube")
		fmt.Println("Starting SonarQube container")
		docker.StartContainer(cli, contId)
	}

	fmt.Println(contId)
}
