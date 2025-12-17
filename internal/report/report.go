package report

import (
	"errors"

	"github.com/sandronister/docker-doctor/internal/checks"
)

type Report struct {
	Results []checks.Result `json:"results"`
	Summary Summary         `json:"summary"`
}

type Summary struct {
	OK   int `json:"ok"`
	WARN int `json:"warn"`
	FAIL int `json:"fail"`
}

type ExitSeverity int

const (
	ExitWarn ExitSeverity = iota
	ExitFail
)

func ParseSeverity(s string) (ExitSeverity, error) {
	switch s {
	case "warn":
		return ExitWarn, nil
	case "fail":
		return ExitFail, nil
	default:
		return ExitWarn, errors.New("invalid severity: use warn|fail")
	}
}

func FromResults(res []checks.Result) Report {
	var sum Summary
	for _, r := range res {
		switch r.Severity {
		case checks.OK:
			sum.OK++
		case checks.WARN:
			sum.WARN++
		case checks.FAIL:
			sum.FAIL++
		}
	}
	return Report{Results: res, Summary: sum}
}

// Exit codes:
// 0 = ok
// 1 = warnings (when severity=warn)
// 2 = failures (always), or warnings when severity=fail
func ExitCode(r Report, sev ExitSeverity) int {
	if r.Summary.FAIL > 0 {
		return 2
	}
	if r.Summary.WARN > 0 {
		if sev == ExitFail {
			return 2
		}
		return 1
	}
	return 0
}
