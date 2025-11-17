package pipeline

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func DockerRunSteps(steps string, containerName string, path string) {
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

	containerConfig := &container.Config{
		Image:      imageName,
		WorkingDir: "/app",
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
		Cmd: []string{"sleep", "60"},
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
		RestartPolicy: container.RestartPolicy{},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/app",
			},
		},
	}
	fmt.Println("Creating container: ", containerName)

	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		fmt.Println("Can't create the container: ", err)
	}
	fmt.Printf("Container created with ID: %s\n", resp.ID)

	err = StartContainer(ctx, cli, resp.ID)
	if err != nil {
		panic(err)
	}

	err = ExecCommand(steps, resp.ID)
	if err != nil {
		panic(err)
	}

	err = StopContainer(ctx, cli, resp.ID)
	if err != nil {
		panic(err)
	}

	err = RemoveContainer(ctx, cli, resp.ID)
	if err != nil {
		panic(err)
	}
}

func StartContainer(ctx context.Context, cli *client.Client, containerID string) error {
	fmt.Printf("Attempting to start the container: %s\n", containerID)
	time.Sleep(1 * time.Second)
	err := cli.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully started the container %s: ", containerID)
	return nil
}

func ExecCommand(command, containerID string) error {
	fmt.Println("On the way of running the commands.....")
	time.Sleep(1 * time.Second)
	cmd := "docker"
	fullCmd := exec.Command(cmd, "exec", "-i", containerID, command)

	fullCmd.Stderr = os.Stderr

	if err := fullCmd.Run(); err != nil {
		return err
	}
	fmt.Println("Running commands compeletion completed.")
	return nil
}

func StopContainer(ctx context.Context, cli *client.Client, containerID string) error {
	fmt.Printf("Attempting to stop the container: %s\n", containerID)

	time.Sleep(1 * time.Second)
	err := cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully stopped the container %s: ", containerID)
	return nil
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
