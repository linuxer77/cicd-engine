package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		http.Error(w, "Invalid Repo link", http.StatusBadRequest)
		return
	}

	if len(p.Steps) == 0 {
		http.Error(w, "Steps should be more than 0", http.StatusBadRequest)
		return
	}
}

func HandleSteps(steps []string) error {
	for i, step := range p.Steps {
		containername := "step" + strconv.Itoa(i)
		pipeline.DockerRunSteps(step, containername)
	}
}
