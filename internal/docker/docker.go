package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"strings"
)

const RunningState = "running"
const shortIdLength = 12

func NewDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return cli
}

func FindContainerByImageName(cli *client.Client, name string) (types.Container, bool) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
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

	return types.Container{}, false
}

func PullImage(cli *client.Client, image string) {
	out, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(io.Discard, out)
}

func CreateContainer(cli *client.Client, image, name string, ports []nat.Port) string {
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	for _, port := range ports {
		exposedPorts[port] = struct{}{}
		portNum := strings.Split(string(port), "/")[0]
		portBindings[port] = []nat.PortBinding{{HostPort: portNum}}
	}

	res, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        image,
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		nil,
		nil,
		name,
	)
	if err != nil {
		panic(err)
	}

	return res.ID
}

func StartContainer(cli *client.Client, id string) {
	err := cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
}

func ShortId(id string) string {
	return id[:shortIdLength]
}

func FollowContainerLogStream(cli *client.Client, id string) io.ReadCloser {
	out, err := cli.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
	})
	if err != nil {
		panic(err)
	}
	return out
}
