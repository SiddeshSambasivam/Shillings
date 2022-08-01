package main

import (
	"log"
	"net"
	"net/http"

	pkg "github.com/SiddeshSambasivam/shillings/pkg"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
)

var PORT = ":8000"
var ADDR = "127.0.0.1"

func loginHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		client, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()

		cmd := &pb.RequestCommand{Command: pb.Command_LGN}
		pkg.SendCmdRequest(client, cmd)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {

	envPort := ":" + pkg.GetEnvVar("WEB_PORT")

	if envPort != "" {
		PORT = envPort
		log.Println("Loaded env var:", envPort)
	}

	ADDR = ADDR + PORT

	http.HandleFunc("/login", loginHandler)

	log.Println("Serving web server @ : " + ADDR)
	err := http.ListenAndServe(ADDR, nil)
	pkg.HandleErrorWithExt(err)

}
