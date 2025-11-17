package pipeline

type Pipeline struct {
	Name  string   `json:"name"`
	Repo  string   `json:"repo"`
	Steps []string `json:"steps"`
}
