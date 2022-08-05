package main

import (
	"log"
	"net"
	"net/http"

	"github.com/SiddeshSambasivam/shillings/pkg/protocols"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		w.Header().Set("Content-Type", "application/json")

		client, err := net.Dial("tcp", "app:8020")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()
		cmd := &pb.RequestCommand{Command: pb.Command_LGN}

		err = protocols.SendProtocolData(client, cmd)
		if err != nil {
			resp := make(map[string]string)

			w.WriteHeader(http.StatusInternalServerError)
			resp["message"] = "Unable to connect to application server"
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
