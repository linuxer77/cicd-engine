package pipeline

type Pipeline struct {
	Name  string   `json:"name"`
	Steps []string `json:"steps"`
}
