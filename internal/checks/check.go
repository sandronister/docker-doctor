package checks

import (
	"context"
	"sync"
	"time"
)

type Severity string

const (
	OK   Severity = "ok"
	WARN Severity = "warn"
	FAIL Severity = "fail"
)

type Result struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Severity    Severity `json:"severity"`
	Summary     string   `json:"summary"`
	Details     string   `json:"details,omitempty"`
	Suggestions []string `json:"suggestion,omitempty"`
	DurationMS  int64    `json:"durationMs"`
}

type Check interface {
	ID() string
	Title() string
	Run(ctx context.Context) Result
}

type RunnerConfig struct {
	Parallelism int
	PerCheckTO  time.Duration
}

type Runner struct {
	cfg RunnerConfig
}

func NewRunner(cfg RunnerConfig) *Runner {
	if cfg.Parallelism <= 0 {
		cfg.Parallelism = 4
	}

	if cfg.PerCheckTO <= 0 {
		cfg.PerCheckTO = 30 * time.Second
	}

	return &Runner{cfg: cfg}
}

func (r *Runner) Run(ctx context.Context, list []Check) []Result {
	type job struct {
		i int
		c Check
	}

	jobs := make(chan job)
	out := make([]Result, len(list))

	var wg sync.WaitGroup
	workers := r.cfg.Parallelism
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				start := time.Now()
				cctx, cancel := context.WithTimeout(ctx, r.cfg.PerCheckTO)
				res := j.c.Run(cctx)
				cancel()
				res.DurationMS = time.Since(start).Milliseconds()
				out[j.i] = res
			}
		}()
	}

	for i, c := range list {
		select {
		case <-ctx.Done():
			break
		case jobs <- job{i: i, c: c}:
		}
	}
	close(jobs)
	wg.Wait()

	return out

}
