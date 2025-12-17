package checks

import (
	"context"
	"fmt"

	"golang.org/x/sys/unix"
)

type diskUsageCheck struct{}

func NewDiskUsageCheck() Check { return &diskUsageCheck{} }

func (c *diskUsageCheck) ID() string    { return "host.disk" }
func (c *diskUsageCheck) Title() string { return "Host disk usage" }

func (c *diskUsageCheck) Run(ctx context.Context) Result {
	var st unix.Statfs_t
	if err := unix.Statfs("/", &st); err != nil {
		return Result{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: WARN,
			Summary:  "Could not read filesystem stats",
			Details:  fmt.Sprintf("%v", err),
		}
	}

	total := float64(st.Blocks) * float64(st.Bsize)
	free := float64(st.Bavail) * float64(st.Bsize)
	used := total - free
	usedPct := (used / total) * 100.0

	sev := OK
	summary := fmt.Sprintf("Used %.1f%% of /", usedPct)

	var sugg []string
	if usedPct >= 95 {
		sev = FAIL
		summary = fmt.Sprintf("CRITICAL: used %.1f%% of /", usedPct)
		sugg = []string{
			"Run `docker system df` to see Docker disk usage",
			"Consider `docker system prune` (careful!)",
			"Clean large files/logs; check /var/lib/docker size",
		}
	} else if usedPct >= 85 {
		sev = WARN
		summary = fmt.Sprintf("High disk usage: %.1f%% of /", usedPct)
		sugg = []string{
			"Investigate large directories; Docker can consume a lot of space",
			"Run `docker system df` to inspect images/volumes",
		}
	}

	return Result{
		ID:          c.ID(),
		Title:       c.Title(),
		Severity:    sev,
		Summary:     summary,
		Suggestions: sugg,
	}
}
