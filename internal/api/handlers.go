package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/linuxer77/cicd/internal/pipeline"
)

var p pipeline.Pipeline

func ParseInst(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(p.Repo)
	if err != nil {
		e := fmt.Sprintln("Error when running the get repo command: ", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != 200 {
		http.Error(w, "Invalid Repo link", http.StatusBadRequest)
		return
	}

	path, err := pipeline.CloneRepo(p.Repo)
	if err != nil {
		e := fmt.Sprintln("Error when running the clone command: ", err)
		http.Error(w, e, http.StatusInternalServerError)
	}

	if len(p.Steps) == 0 {
		http.Error(w, "Steps should be more than 0", http.StatusBadRequest)
		return
	}

	err = HandleSteps(p.Steps, path)
	if err != nil {
		http.Error(w, "some error has occurred when handling steps", http.StatusInternalServerError)
		fmt.Println("Error when handling steps: ", err)
		return
	}

	pipeline.RemoveDir(p.Repo)
}

func HandleSteps(steps []string, path string) error {
	for i, step := range p.Steps {
		command := strings.Fields(step)
		containername := "step" + strconv.Itoa(i)
		pipeline.DockerRunSteps(command, containername, path)
	}
	return nil
}
