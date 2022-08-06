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

	http.HandleFunc("/v1/login", loginHandler)
	http.HandleFunc("/v1/signup", signupHandler)
	http.HandleFunc("/v1/account", userAccountHandler)
	http.HandleFunc("/v1/pay", paymentHandler)

	log.Println("Serving web server @ : " + ADDR)
	err := http.ListenAndServe(ADDR, nil)
	errors.HandleErrorWithExt(err)

}
