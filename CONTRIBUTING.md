
---

## GitHub Actions: .github/workflows/ci.yml

```yaml
name: CI

on:
  push:
  pull_request:

jobs:
  build-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Test
        run: go test ./... -v
      - name: Build
        run: go build ./cmd/dockdoc
