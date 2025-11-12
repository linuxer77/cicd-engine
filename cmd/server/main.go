package main

import (
	"github.com/linuxer77/cicd/internal/pipeline"
)

func main() {
	pipeline.CloneRepo("https://github.com/linuxer77/trade-execution")
}
