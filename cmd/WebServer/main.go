package main

import (
	"log"
	"net"
	"net/http"

	envtools "github.com/SiddeshSambasivam/shillings/pkg/env"
	errors "github.com/SiddeshSambasivam/shillings/pkg/errors"
)

var p Pool

func main() {

	var PORT = ":8000"
	var ADDR = "0.0.0.0"
	envPort := ":" + envtools.GetEnvVar("WEB_PORT")

	p, _ = NewChannelPool(10, 20000, func() (net.Conn, error) { return net.Dial("tcp", "app:8020") })

	if envPort != "" {
		PORT = envPort
		log.Println("Loaded env var:", envPort)
	}

	ADDR = ADDR + PORT

	http.HandleFunc("/v1/login", loginHandler)
	http.HandleFunc("/v1/signup", signupHandler)
	http.HandleFunc("/v1/account", userAccountHandler)
	http.HandleFunc("/v1/pay", paymentHandler)
	http.HandleFunc("/v1/topup", topupHandler)
	http.HandleFunc("/v1/transactions", transactionsHandler)

	log.Println("Serving web server @ : " + ADDR)
	err := http.ListenAndServe(ADDR, nil)
	errors.HandleErrorWithExt(err)

}
