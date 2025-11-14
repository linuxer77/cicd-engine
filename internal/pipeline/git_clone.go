package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CloneRepo(url string) {
	parts := strings.Split(url, "/")
	fmt.Println("parts: ", parts)
	location := fmt.Sprintf("/tmp/repos/%s", parts[len(parts)-1])
	fmt.Println("location: ", location)

	cmd := exec.Command("git", "clone", url, location)

	cmd.Stdout = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command: ", err)
	}
}
