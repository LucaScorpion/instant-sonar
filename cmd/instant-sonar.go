package main

import (
	"fmt"
	"instant-sonar/internal"
)

func main() {
	fmt.Println("Starting instant Sonar")
	cli := internal.NewDockerClient()
	defer cli.Close()

	contId := ""
	if cont, exists := internal.FindContainerByImageName(cli, internal.SonarqubeImage); exists {
		fmt.Println("SonarQube container already exists")
		contId = cont.ID

		if cont.State != "running" {
			fmt.Println("Starting SonarQube container")
			internal.StartContainer(cli, contId)
		}
	} else {
		fmt.Println("Pulling SonarQube image")
		internal.PullImage(cli, internal.SonarqubeImage)
		fmt.Println("Creating SonarQube container")
		contId = internal.CreateContainer(cli, internal.SonarqubeImage, "sonarqube")
		fmt.Println("Starting SonarQube container")
		internal.StartContainer(cli, contId)
	}

	fmt.Println(contId)
}
