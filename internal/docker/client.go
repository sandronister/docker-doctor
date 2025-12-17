package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type Client interface {
	Info(ctx context.Context) (any, error)
	ContainerList(ctx context.Context, options any) ([]any, error)
	ContainerInspect(ctx context.Context, containerID string) (any, error)
	Ping(ctx context.Context, options client.PingOptions) (any, error)
}

func NewClient() (*client.Client, error) {
	return client.New(client.FromEnv)
}
