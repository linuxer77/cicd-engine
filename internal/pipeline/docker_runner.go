package pipeline

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func DockerRunSteps() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	imageName := "golang:1.23"
	fmt.Printf("Pulling your image: %s\n", imageName)

	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
	reader.Close()
	fmt.Println("Successfully pulled the image")

	containerName := "step2"

	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"80/tcp": struct{}{},
		},
		Env: []string{
			"NGINX_HOST=localhost",
			"NGINX_PORT=80",
		},
		Labels: map[string]string{
			"created-by": "my-go-program",
			"purpose":    "testing",
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}
	fmt.Println("Creating container: ", containerName)

	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil, // network config
		nil, // platform
		containerName,
	)
	if err != nil {
		fmt.Println("Can't create the container: ", err)
	}
	fmt.Println("resp: ", resp)
	fmt.Printf("Container created with ID: %s\n", resp.ID)
	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Container started Successfully")

	fmt.Println("On the way of running the commands.....")

	fmt.Println("Stopping container now.")
	err = StopContainer(ctx, cli, resp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully stopped the container")

	err = RemoveContainer(ctx, cli, resp.ID)
	if err != nil {
		panic(err)
	}
}

func RemoveContainer(ctx context.Context, cli *client.Client, containerID string) error {
	fmt.Printf("Attempting to remove the container: %s\n", containerID)
	err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully removed the container: %s\n", containerID)
	return nil
}

func StopContainer(ctx context.Context, cli *client.Client, containerID string) error {
	fmt.Printf("Attempting to stop the container: %s\n", containerID)

	err := cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully stopped the container %s: ", containerID)
	return nil
}

func ExecCommands(step string) error {
	formattedCmds := strings.Fields(step)
	name := formattedCmds[0]
	cmd := exec.Command(name, formattedCmds[1:]...)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command")
		return err
	}
	return nil
}
