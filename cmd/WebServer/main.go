package main

import (
	"log"
	"net/http"

	envtools "github.com/SiddeshSambasivam/shillings/pkg/env"
	errors "github.com/SiddeshSambasivam/shillings/pkg/errors"
)

func main() {

	var PORT = ":8000"
	var ADDR = "0.0.0.0"
	envPort := ":" + envtools.GetEnvVar("WEB_PORT")

	if envPort != "" {
		PORT = envPort
		log.Println("Loaded env var:", envPort)
	}

	ADDR = ADDR + PORT

	http.HandleFunc("/login", loginHandler)

	log.Println("Serving web server @ : " + ADDR)
	err := http.ListenAndServe(ADDR, nil)
	errors.HandleErrorWithExt(err)

}
