package runner

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TestDockerNetworkList(t *testing.T) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Fatalf("new env client: %s", err)
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		t.Fatalf("network list: %s", err)
	}

	fmt.Println("networks:")
	for _, network := range networks {
		fmt.Printf("%s %s\n", network.ID[:10], network.Name)
	}
}

func TestDockerContainerList(t *testing.T) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Fatalf("new env client: %s", err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		t.Fatalf("container list: %s", err)
	}

	fmt.Println("containers:")
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}
