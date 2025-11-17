package pipeline

import (
	"fmt"
	"os/exec"
	"strings"
)

type Ouptut struct {
	op string `json:"op"`
}

func RunCmds(steps []string) error {
	for _, step := range steps {
		formattedCmds := strings.Fields(step)
		name := formattedCmds[0]
		cmd := exec.Command(name, formattedCmds[1:]...)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("could not run command")
			return err
		}
		fmt.Printf("output: %s\n", string(out))
	}
	return nil
}
