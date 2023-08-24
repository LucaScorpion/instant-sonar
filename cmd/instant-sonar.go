package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	flag "github.com/spf13/pflag"
	"instant-sonar/internal"
	"instant-sonar/internal/docker"
	"instant-sonar/internal/log"
	"instant-sonar/internal/sonar"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type options struct {
	help    bool
	verbose bool

	username string
	password string

	path string
	copy bool
}

var flags *flag.FlagSet

func initCli() (*options, error) {
	flags = flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	opts := options{}

	flags.BoolVarP(&opts.help, "help", "h", false, "Print info and help")
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "More verbose logging")
	flags.StringVarP(&opts.username, "username", "u", "admin", "SonarQube admin username")
	flags.StringVarP(&opts.password, "password", "p", "admin", "SonarQube admin password")
	flags.BoolVarP(&opts.copy, "copy", "c", false, "Copy the files into the Sonar Scanner container instead of using a bound volume\nThis is generally faster on Mac and Windows")

	err := flags.Parse(os.Args[1:])
	opts.path, _ = filepath.Abs(flags.Arg(0))

	return &opts, err
}

func main() {
	opts, err := initCli()
	log.IsVerbose = opts.verbose

	if flags.NArg() > 1 {
		err = errors.New("too many paths given")
	}

	if err != nil && !errors.Is(err, flag.ErrHelp) {
		fmt.Println(err)
		printHelp()
		os.Exit(2)
	}

	if opts.help {
		printHelp()
		return
	}

	log.Verboseln("Starting Instant Sonar")
	cli := docker.NewClient()
	defer cli.Close()

	log.Println("Preparing SonarQube")

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
	out := cli.FollowContainerLogStream(qubeContId, cli.InspectContainer(qubeContId).State.StartedAt)
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

	sonarApi := sonar.NewApiClient("http://127.0.0.1:9000", opts.username, opts.password)

	log.Verboseln("Checking SonarQube API")
	if err := sonarApi.Ping(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	log.Verboseln("Disabling force user auth")
	sonarApi.DisableForceUserAuth()

	projectKey := internal.RandomString(16)
	log.Verboseln("Creating project (" + projectKey + ")")
	sonarApi.CreateProject(filepath.Base(opts.path), projectKey)

	log.Verbose("Creating analysis token")
	token := sonarApi.CreateToken(projectKey)
	log.Verboseln(" (" + token + ")")

	log.Println("Preparing Sonar Scanner")

	log.Verboseln("Pulling Sonar Scanner image")
	cli.PullImage(sonar.SonarScannerImage)

	log.Verbose("Creating Sonar Scanner container")
	qubeDockerUrl := "http://" + cli.GetContainerIp(qubeContId) + ":9000"
	curUser, _ := user.Current()
	scanContId := sonar.CreateSonarScannerContainer(cli, qubeDockerUrl, projectKey, token, opts.path, curUser.Uid, !opts.copy)
	log.Verboseln(" (" + docker.ShortId(scanContId) + ")")

	if opts.copy {
		log.Verboseln("Copying files to Sonar Scanner container")
		cli.CopyDirToContainer(scanContId, opts.path, "/usr/src")
	}

	log.Println("Starting analysis")
	cli.StartContainer(scanContId)
	cli.WaitForContainer(scanContId, container.WaitConditionRemoved)

	log.Verboseln("Removing scannerwork directory")
	os.RemoveAll(filepath.Join(opts.path, sonar.ScannerworkDir))

	log.Println("Project dashboard: " + sonarApi.ProjectDashboardUrl(projectKey))
}

func printHelp() {
	cmdName := filepath.Base(os.Args[0])
	fmt.Println("Usage: " + cmdName + " [options...] [path]")
	flags.PrintDefaults()
}
