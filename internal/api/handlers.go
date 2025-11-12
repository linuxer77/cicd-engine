package api

import (
	"encoding/json"
	"log"
	"net/http"

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

	err = pipeline.RunCmds(p.Steps)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}
