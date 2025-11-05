package api

import (
	"encoding/json"
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

	if len(p.Steps) == 0 {
		http.Error(w, "Steps should be more than 0", http.StatusBadRequest)
		return
	}

	err = pipeline.RunCmds()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}
