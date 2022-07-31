package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var PORT = "8080"
var ADDR = "127.0.0.1"

func getEnvVar(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)

	resp["message"] = ADDR
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)

}

func main() {

	envPort := getEnvVar("PORT")
	log.Println("Loaded env var: ", envPort)

	if envPort != "" {
		PORT = envPort
	}

	ADDR = ADDR + ":" + PORT

	log.Println("Starting server on " + ADDR)
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(ADDR, nil)

	if err != nil {
		log.Fatal(err)
	}
}
