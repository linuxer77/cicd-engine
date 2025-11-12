package pipeline

type Pipeline struct {
	Name  string   `json:"name"`
	Repo  string   `json:"repo"`
	Image string   `json:"image"`
	Steps []string `json:"steps"`
}
