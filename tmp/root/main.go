package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://github.com/linuxer77/VeriCred")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.StatusCode)
}
