package checks

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

type dockerDaemonCheck struct {
	cli *client.Client
	err error
}

func NewDockerDaemonCheck(cli *client.Client, initErr error) Check {
	return &dockerDaemonCheck{cli: cli, err: initErr}
}

func (c *dockerDaemonCheck) ID() string    { return "docker.daemon" }
func (c *dockerDaemonCheck) Title() string { return "Docker daemon reachable" }

func (c *dockerDaemonCheck) Run(ctx context.Context) Result {
	if c.err != nil || c.cli == nil {
		return Result{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: FAIL,
			Summary:  "Docker client could not be created",
			Details:  fmt.Sprintf("%v", c.err),
			Suggestions: []string{
				"On macOS/Windows: ensure Docker Desktop is running",
				"On Linux: check `systemctl status docker` and user permissions (docker group)",
				"Verify DOCKER_HOST / context settings",
			},
		}
	}

	p, err := c.cli.Ping(ctx, client.PingOptions{})
	if err != nil {
		return Result{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: FAIL,
			Summary:  "Docker daemon not responding to ping",
			Details:  fmt.Sprintf("%v", err),
			Suggestions: []string{
				"Restart Docker (Desktop) or `sudo systemctl restart docker`",
				"Check if the Docker socket is accessible (permissions)",
			},
		}
	}

	return Result{
		ID:       c.ID(),
		Title:    c.Title(),
		Severity: OK,
		Summary:  fmt.Sprintf("OK (API %s)", p.APIVersion),
	}
}
