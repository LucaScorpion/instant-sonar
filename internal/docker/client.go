package docker

import (
	"context"
	dockerTypes "github.com/docker/docker/api/types"
	dockerSdk "github.com/docker/docker/client"
	"io"
	"strings"
)

const RunningState = "running"
const shortIdLength = 12

type Client struct {
	Cli *dockerSdk.Client
}

func NewClient() *Client {
	cli, err := dockerSdk.NewClientWithOpts(dockerSdk.FromEnv, dockerSdk.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &Client{Cli: cli}
}

func (c *Client) Close() {
	c.Cli.Close()
}

func (c *Client) FindContainerByImageName(name string) (dockerTypes.Container, bool) {
	containers, err := c.Cli.ContainerList(context.Background(), dockerTypes.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	for _, cont := range containers {
		imageName := strings.Split(cont.Image, ":")[0]
		if imageName == name {
			return cont, true
		}
	}

	return dockerTypes.Container{}, false
}

func (c *Client) PullImage(image string) {
	out, err := c.Cli.ImagePull(context.Background(), image, dockerTypes.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(io.Discard, out)
}

func (c *Client) StartContainer(id string) {
	err := c.Cli.ContainerStart(context.Background(), id, dockerTypes.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
}

func (c *Client) FollowContainerLogStream(id string) io.ReadCloser {
	out, err := c.Cli.ContainerLogs(context.Background(), id, dockerTypes.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
	})
	if err != nil {
		panic(err)
	}
	return out
}

func (c *Client) GetContainerIp(id string) string {
	return "http://127.0.0.1:9000" // TODO
}