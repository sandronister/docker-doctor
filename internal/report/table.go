package report

import (
	"fmt"
	"strings"
)

func RenderTable(r Report) string {
	var b strings.Builder

	b.WriteString("DOCKER DOCTOR REPORT\n")
	b.WriteString(fmt.Sprintf("Summary: OK=%d WARN=%d FAIL=%d\n\n", r.Summary.OK, r.Summary.WARN, r.Summary.FAIL))

	for _, it := range r.Results {
		b.WriteString(fmt.Sprintf("[%s] %s\n", strings.ToUpper(string(it.Severity)), it.Title))
		b.WriteString(fmt.Sprintf("  - %s\n", it.Summary))
		if it.Details != "" {
			for _, line := range strings.Split(it.Details, "\n") {
				b.WriteString(fmt.Sprintf("  - %s\n", line))
			}
		}
		if len(it.Suggestions) > 0 {
			b.WriteString("  Suggestions:\n")
			for _, s := range it.Suggestions {
				b.WriteString(fmt.Sprintf("    * %s\n", s))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}
