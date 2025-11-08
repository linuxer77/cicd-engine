package pipeline

import (
	"fmt"
	"os"
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
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println("could not run command")
			return err
		}
	}
	return nil
}
