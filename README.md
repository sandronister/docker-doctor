# docker-doctor (dockdoc)

Fast diagnostics for broken Docker environments — with actionable fixes.

## Why this exists
When Docker breaks, people lose hours guessing: daemon down, disk full, containers restarting, permissions, networking...
`dockdoc` runs focused checks and prints a clear report.

## What it does
- Checks Docker daemon connectivity (ping)
- Checks host disk usage (common Docker failure cause)
- Detects containers stuck restarting
- Outputs `table` or `json`, with CI-friendly exit codes

## What it does NOT do
- It does not “auto-fix” your machine by default
- It is not a Kubernetes tool
- It is not a Docker management UI

## Install
### From source
```bash
go install github.com/sandronister/docker-doctor/cmd/dockdoc@latest
```

## Quick Start

```bash 
dockdoc scan
dockdoc scan --format json --severity fail 
```
