package main

import (
	"bufio"
	"github.com/docker/docker/api/types/container"
	flag "github.com/spf13/pflag"
	"instant-sonar/internal"
	"instant-sonar/internal/docker"
	"instant-sonar/internal/log"
	"instant-sonar/internal/sonar"
	"os"
	"os/user"
	"path"
	"strings"
)

type options struct {
	verbose bool
}

func initCli() *options {
	opts := options{}
	flag.BoolVarP(&opts.verbose, "verbose", "v", false, "More verbose logging")
	flag.Parse()
	return &opts
}

func main() {
	opts := initCli()
	log.IsVerbose = opts.verbose

	log.Println("Starting Instant Sonar")

	cli := docker.NewClient()
	defer cli.Close()

	log.Println("Preparing SonarQube image")
	qubeContId := ""
	if cont, exists := cli.FindContainerByImageName(sonar.SonarqubeImage); exists {
		qubeContId = cont.ID
		log.Verboseln("SonarQube container already exists (" + docker.ShortId(qubeContId) + ")")

		if cont.State != docker.RunningState {
			log.Verboseln("Starting SonarQube container")
			cli.StartContainer(qubeContId)
		}
	} else {
		log.Verboseln("Pulling SonarQube image")
		cli.PullImage(sonar.SonarqubeImage)

		log.Verbose("Creating SonarQube container")
		qubeContId = sonar.CreateSonarQubeContainer(cli)
		log.Verboseln(" (" + docker.ShortId(qubeContId) + ")")

		log.Verboseln("Starting SonarQube container")
		cli.StartContainer(qubeContId)
	}

	log.Verboseln("Waiting for SonarQube to be operational")
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
	log.Verboseln("SonarQube is operational")

	log.Println("Preparing Sonar Scanner image")
	qubeWebUrl := "http://" + cli.GetContainerIp(qubeContId) + ":9000"
	sonarApi := sonar.NewApiClient(qubeWebUrl, "admin", "admin")
	sonarApi.DisableForceUserAuth()

	log.Verbose("Creating project")
	projectKey := internal.RandomString(16)
	sonarApi.CreateProject(projectKey, projectKey)
	log.Verboseln(" (" + projectKey + ")")

	log.Verbose("Creating analysis token")
	token := sonarApi.CreateToken(projectKey)
	log.Verboseln(" (" + token + ")")

	log.Verboseln("Pulling Sonar Scanner image")
	cli.PullImage(sonar.SonarScannerImage)

	log.Verbose("Creating Sonar Scanner container")
	scanDir, _ := os.Getwd() // TODO: Get from argument
	curUser, _ := user.Current()
	scanContId := sonar.CreateSonarScannerContainer(cli, sonarApi.Url, projectKey, token, scanDir, curUser.Uid)
	log.Verboseln(" (" + docker.ShortId(scanContId) + ")")

	log.Println("Starting analysis")
	cli.StartContainer(scanContId)
	cli.WaitForContainer(scanContId, container.WaitConditionRemoved)

	log.Verboseln("Removing scannerwork directory")
	os.RemoveAll(path.Join(scanDir, sonar.ScannerworkDir))

	log.Println("Project dashboard: " + sonarApi.ProjectDashboardUrl(projectKey))
}
