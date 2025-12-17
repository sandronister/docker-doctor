package checks

import (
	"context"
	"fmt"
	"strings"

	"github.com/moby/moby/client"
)

type restartingContainersCheck struct {
	cli *client.Client
}

func NewRestartingContainersCheck(cli *client.Client) Check {
	return &restartingContainersCheck{cli: cli}
}

func (c *restartingContainersCheck) ID() string    { return "docker.restarts" }
func (c *restartingContainersCheck) Title() string { return "Restarting containers" }

func (c *restartingContainersCheck) Run(ctx context.Context) Result {
	if c.cli == nil {
		return Result{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: WARN,
			Summary:  "Skipped (Docker client unavailable)",
		}
	}

	cts, err := c.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return Result{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: WARN,
			Summary:  "Could not list containers",
			Details:  fmt.Sprintf("%v", err),
		}
	}

	var bad []string
	for _, ct := range cts.Items {
		// heurística simples: status contém "Restarting" ou "restarting"
		if strings.Contains(strings.ToLower(ct.Status), "restarting") {
			name := ""
			if len(ct.Names) > 0 {
				name = strings.TrimPrefix(ct.Names[0], "/")
			} else {
				name = ct.ID[:12]
			}
			bad = append(bad, fmt.Sprintf("%s (%s)", name, ct.Status))
		}
	}

	if len(bad) == 0 {
		return Result{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: OK,
			Summary:  "No restarting containers detected",
		}
	}

	return Result{
		ID:       c.ID(),
		Title:    c.Title(),
		Severity: WARN,
		Summary:  fmt.Sprintf("%d container(s) restarting", len(bad)),
		Details:  strings.Join(bad, "\n"),
		Suggestions: []string{
			"Check logs: `docker logs <container>`",
			"Inspect config: `docker inspect <container>`",
			"Common causes: missing env vars, ports in use, DB unavailable",
		},
	}
}
