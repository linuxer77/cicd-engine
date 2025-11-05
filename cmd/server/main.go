package main

import (
	"net/http"

	"github.com/linuxer77/cicd/internal/api"
)

func main() {
	r := api.Router()
	http.ListenAndServe(":8080", r)
}
