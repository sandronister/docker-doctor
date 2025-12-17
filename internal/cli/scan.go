package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sandronister/docker-doctor/internal/checks"
	"github.com/sandronister/docker-doctor/internal/docker"
	"github.com/sandronister/docker-doctor/internal/report"
	"github.com/spf13/cobra"
)

var (
	scanFormat   string
	scanSeverity string
	scanTimeout  time.Duration
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run diagnostics checks and print a report",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), scanTimeout)
		defer cancel()

		sev, err := report.ParseSeverity(scanSeverity)
		if err != nil {
			return err
		}

		cli, derr := docker.NewClient()
		if derr != nil {
			// Still run checks that don't need Docker API.
			cli = nil
		}

		runner := checks.NewRunner(checks.RunnerConfig{
			Parallelism: 4,
			PerCheckTO:  8 * time.Second,
		})

		results := runner.Run(ctx, []checks.Check{
			checks.NewDockerDaemonCheck(cli, derr),
			checks.NewDiskUsageCheck(),
			checks.NewRestartingContainersCheck(cli),
		})

		rep := report.FromResults(results)

		switch scanFormat {
		case "table":
			fmt.Print(report.RenderTable(rep))
		case "json":
			b, jerr := report.RenderJSON(rep)
			if jerr != nil {
				return jerr
			}
			fmt.Println(string(b))
		default:
			return errors.New("invalid --format, use table|json")
		}

		code := report.ExitCode(rep, sev)
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}

func init() {
	scanCmd.Flags().StringVar(&scanFormat, "format", "table", "Output format: table|json")
	scanCmd.Flags().StringVar(&scanSeverity, "severity", "warn", "Exit level: warn|fail")
	scanCmd.Flags().DurationVar(&scanTimeout, "timeout", 30*time.Second, "Total scan timeout")
}
