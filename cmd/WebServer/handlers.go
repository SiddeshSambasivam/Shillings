package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/SiddeshSambasivam/shillings/pkg/models"
	"github.com/SiddeshSambasivam/shillings/pkg/protocols"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"google.golang.org/protobuf/proto"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		resp := make(map[string]string)
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user models.User
		var credentials models.Credentials

		json.Unmarshal(data, &user)
		json.Unmarshal(data, &credentials)

		request := &pb.RequestSignup{
			User: &pb.User{
				FirstName:  user.First_name,
				MiddleName: user.Middle_name,
				LastName:   user.Last_name,
				Email:      user.Email,
				Phone:      user.Phone,
				Balance:    0.0,
			},
			Credentials: &pb.Credentials{
				Email:    credentials.Email,
				Password: credentials.Password,
			},
		}

		client, err := net.Dial("tcp", "app:8020")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()
		cmd := &pb.RequestCommand{Command: pb.Command_SGN}

		err = protocols.SendProtocolData(client, cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["message"] = err.Error()
		}

		err = protocols.SendProtocolData(client, request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["message"] = err.Error()
		}

		response := &pb.ResponseSignup{}
		respBytes, err := protocols.ReadProtocolData(client)
		if err != nil {
			log.Println("Error reading data from application server: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = proto.Unmarshal(respBytes, response)
		if err != nil {
			log.Println("Error unmarshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		jsonData := SignUpResponse{
			Status: Status{
				Code:    int32(response.GetStatus().GetCode()),
				Message: response.GetStatus().GetMessage(),
			},
		}

		jsonResp, err := json.Marshal(jsonData)
		if err != nil {
			log.Println("Error marshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if jsonData.Status.Code == int32(pb.Code_OK) {
			w.WriteHeader(http.StatusOK)
		}
		w.Write(jsonResp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

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
