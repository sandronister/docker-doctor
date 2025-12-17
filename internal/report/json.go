package report

import "encoding/json"

func RenderJSON(r Report) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
