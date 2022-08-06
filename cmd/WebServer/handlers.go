package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/SiddeshSambasivam/shillings/pkg/models"
	"github.com/SiddeshSambasivam/shillings/pkg/protocols"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"google.golang.org/protobuf/proto"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

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
			return
		}

		err = protocols.SendProtocolData(client, request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &pb.ResponseSignup{}
		respBytes, err := protocols.ReadProtocolData(client)
		if err != nil {
			log.Println("Error reading data from application server: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = proto.Unmarshal(respBytes, response)
		if err != nil {
			log.Println("Error unmarshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData := SignUpResponse{
			Status: Status{
				Code:    int32(response.GetStatus().GetCode()),
				Message: response.GetStatus().GetMessage(),
			},
		}

		switch jsonData.Status.Code {
		case int32(pb.Code_OK):
			w.WriteHeader(http.StatusOK)
		case int32(pb.Code_BAD_REQUEST):
			w.WriteHeader(http.StatusBadRequest)
		case int32(pb.Code_INTERNAL_SERVER_ERROR):
			w.WriteHeader(http.StatusInternalServerError)
		case int32(pb.Code_FORBIDDEN):
			w.WriteHeader(http.StatusForbidden)
		case int32(pb.Code_Conflict):
			w.WriteHeader(http.StatusConflict)
		}

		jsonResp, err := json.Marshal(jsonData)
		if err != nil {
			log.Println("Error marshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
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

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var credentials models.Credentials

		json.Unmarshal(data, &credentials)

		client, err := net.Dial("tcp", "app:8020")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()
		cmd := &pb.RequestCommand{Command: pb.Command_LGN}

		err = protocols.SendProtocolData(client, cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		request := &pb.RequestLogin{
			Credentials: &pb.UserCredentials{
				Email:    credentials.Email,
				Password: credentials.Password,
			},
		}

		err = protocols.SendProtocolData(client, request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &pb.ResponseLogin{}
		respBytes, err := protocols.ReadProtocolData(client)
		if err != nil {
			log.Println("Error reading data from application server: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = proto.Unmarshal(respBytes, response)
		if err != nil {
			log.Println("Error unmarshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonData := LoginResponse{
			Status: Status{
				Code:    int32(response.GetStatus().GetCode()),
				Message: response.GetStatus().GetMessage(),
			},
		}

		if jsonData.Status.Code == int32(pb.Code_OK) {
			// set the cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   response.GetAuth().GetToken(),
				Expires: time.Unix(response.GetAuth().GetExpirationTime(), 0),
			})
		}

		switch jsonData.Status.Code {
		case int32(pb.Code_UNAUTHORIZED):
			w.WriteHeader(http.StatusUnauthorized)
		case int32(pb.Code_OK):
			w.WriteHeader(http.StatusOK)
		case int32(pb.Code_BAD_REQUEST):
			w.WriteHeader(http.StatusBadRequest)
		case int32(pb.Code_INTERNAL_SERVER_ERROR):
			w.WriteHeader(http.StatusInternalServerError)
			return
		case int32(pb.Code_FORBIDDEN):
			w.WriteHeader(http.StatusForbidden)
		case int32(pb.Code_Conflict):
			w.WriteHeader(http.StatusConflict)
		}

		jsonResp, err := json.Marshal(jsonData)
		if err != nil {
			log.Println("Error marshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(jsonResp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		w.Header().Set("Content-Type", "application/json")
		authToken, _ := r.Cookie("token")
		if authToken == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var paymentReq PaymentRequest

		json.Unmarshal(data, &paymentReq)
		if paymentReq.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		client, err := net.Dial("tcp", "app:8020")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()
		cmd := &pb.RequestCommand{Command: pb.Command_PAY}

		err = protocols.SendProtocolData(client, cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		request := &pb.RequestPayUser{
			Auth: &pb.Auth{
				Token: authToken.Value,
			},
			ReceiverEmail: paymentReq.Receiver_email,
			Amount:        paymentReq.Amount,
		}

		err = protocols.SendProtocolData(client, request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &pb.ResponsePayUser{}
		respBytes, err := protocols.ReadProtocolData(client)
		if err != nil {
			log.Println("Error reading data from application server: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = proto.Unmarshal(respBytes, response)
		if err != nil {
			log.Println("Error unmarshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData := PaymentResponse{
			Status: Status{
				Code:    int32(response.GetStatus().GetCode()),
				Message: response.GetStatus().GetMessage(),
			},
			Transaction_id: response.GetTransactionId(),
		}

		jsonResp, err := json.Marshal(jsonData)
		if err != nil {
			log.Println("Error marshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch jsonData.Status.Code {
		case int32(pb.Code_OK):
			w.WriteHeader(http.StatusOK)
		case int32(pb.Code_UNAUTHORIZED):
			w.WriteHeader(http.StatusUnauthorized)
		case int32(pb.Code_BAD_REQUEST):
			w.WriteHeader(http.StatusBadRequest)
		case int32(pb.Code_INTERNAL_SERVER_ERROR):
			w.WriteHeader(http.StatusInternalServerError)
			return
		case int32(pb.Code_FORBIDDEN):
			w.WriteHeader(http.StatusForbidden)
		case int32(pb.Code_Conflict):
			w.WriteHeader(http.StatusConflict)
		}

		w.Write(jsonResp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func userAccountHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")

		authToken, _ := r.Cookie("token")
		if authToken == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		client, err := net.Dial("tcp", "app:8020")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()
		cmd := &pb.RequestCommand{Command: pb.Command_USR}

		err = protocols.SendProtocolData(client, cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		request := &pb.RequestGetUser{
			Auth: &pb.Auth{
				Token: authToken.Value,
			},
		}
		err = protocols.SendProtocolData(client, request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &pb.ResponseGetUser{}
		respBytes, err := protocols.ReadProtocolData(client)
		if err != nil {
			log.Println("Error reading data from application server: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = proto.Unmarshal(respBytes, response)
		if err != nil {
			log.Println("Error unmarshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData := UserResponse{
			User: models.User{
				User_id:     response.GetUser().GetUserId(),
				First_name:  response.GetUser().GetFirstName(),
				Middle_name: response.GetUser().GetMiddleName(),
				Last_name:   response.GetUser().GetLastName(),
				Email:       response.GetUser().GetEmail(),
				Phone:       response.GetUser().GetPhone(),
				Balance:     response.GetUser().GetBalance(),
				Created_at:  response.GetUser().GetCreatedAt(),
				Updated_at:  response.GetUser().GetUpdatedAt(),
			},
			Status: Status{
				Code:    int32(response.GetStatus().GetCode()),
				Message: response.GetStatus().GetMessage(),
			},
		}

		jsonResp, err := json.Marshal(jsonData)
		if err != nil {
			log.Println("Error marshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch jsonData.Status.Code {
		case int32(pb.Code_OK):
			w.WriteHeader(http.StatusOK)
		case int32(pb.Code_UNAUTHORIZED):
			w.WriteHeader(http.StatusUnauthorized)
		case int32(pb.Code_BAD_REQUEST):
			w.WriteHeader(http.StatusBadRequest)
		case int32(pb.Code_INTERNAL_SERVER_ERROR):
			w.WriteHeader(http.StatusInternalServerError)
			return
		case int32(pb.Code_FORBIDDEN):
			w.WriteHeader(http.StatusForbidden)
		case int32(pb.Code_Conflict):
			w.WriteHeader(http.StatusConflict)
		}

		w.Write(jsonResp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func topupHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		authToken, _ := r.Cookie("token")
		if authToken == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var topupreq TopupRequest

		json.Unmarshal(data, &topupreq)
		if topupreq.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		client, err := net.Dial("tcp", "app:8020")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()
		cmd := &pb.RequestCommand{Command: pb.Command_TPU}

		err = protocols.SendProtocolData(client, cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		request := &pb.RequestTopupUser{
			Auth: &pb.Auth{
				Token: authToken.Value,
			},
			Amount: topupreq.Amount,
		}

		err = protocols.SendProtocolData(client, request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &pb.ResponseTopupUser{}
		respBytes, err := protocols.ReadProtocolData(client)
		if err != nil {
			log.Println("Error reading data from application server: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = proto.Unmarshal(respBytes, response)
		if err != nil {
			log.Println("Error unmarshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData := TopupResponse{
			Status: Status{
				Code:    int32(response.GetStatus().GetCode()),
				Message: response.GetStatus().GetMessage(),
			},
		}

		jsonResp, err := json.Marshal(jsonData)
		if err != nil {
			log.Println("Error marshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch jsonData.Status.Code {
		case int32(pb.Code_OK):
			w.WriteHeader(http.StatusOK)
		case int32(pb.Code_UNAUTHORIZED):
			w.WriteHeader(http.StatusUnauthorized)
		case int32(pb.Code_BAD_REQUEST):
			w.WriteHeader(http.StatusBadRequest)
		case int32(pb.Code_INTERNAL_SERVER_ERROR):
			w.WriteHeader(http.StatusInternalServerError)
			return
		case int32(pb.Code_FORBIDDEN):
			w.WriteHeader(http.StatusForbidden)
		case int32(pb.Code_Conflict):
			w.WriteHeader(http.StatusConflict)
		}

		w.Write(jsonResp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")

		authToken, _ := r.Cookie("token")
		if authToken == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		client, err := net.Dial("tcp", "app:8020")
		if err != nil {
			log.Println("Error dialing:", err)
			return
		}

		defer client.Close()
		cmd := &pb.RequestCommand{Command: pb.Command_TXQ}

		err = protocols.SendProtocolData(client, cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		request := &pb.RequestGetUser{
			Auth: &pb.Auth{
				Token: authToken.Value,
			},
		}
		err = protocols.SendProtocolData(client, request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &pb.ResponseGetUserTransactions{}
		respBytes, err := protocols.ReadProtocolData(client)
		if err != nil {
			log.Println("Error reading data from application server: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = proto.Unmarshal(respBytes, response)
		if err != nil {
			log.Println("Error unmarshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var transactions []Transaction = []Transaction{}
		for _, tx := range response.GetTransactions() {
			transactions = append(transactions, Transaction{
				Sender_email:   tx.GetSenderEmail(),
				Receiver_email: tx.GetReceiverEmail(),
				Amount:         tx.GetAmount(),
				Created_at:     tx.GetCreatedAt(),
			})
		}

		jsonData := TransactionResponse{
			Transactions: transactions,
			Status: Status{
				Code:    int32(response.GetStatus().GetCode()),
				Message: response.GetStatus().GetMessage(),
			},
		}

		jsonResp, err := json.Marshal(jsonData)
		if err != nil {
			log.Println("Error marshalling response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch jsonData.Status.Code {
		case int32(pb.Code_OK):
			w.WriteHeader(http.StatusOK)
		case int32(pb.Code_UNAUTHORIZED):
			w.WriteHeader(http.StatusUnauthorized)
		case int32(pb.Code_BAD_REQUEST):
			w.WriteHeader(http.StatusBadRequest)
		case int32(pb.Code_INTERNAL_SERVER_ERROR):
			w.WriteHeader(http.StatusInternalServerError)
			return
		case int32(pb.Code_FORBIDDEN):
			w.WriteHeader(http.StatusForbidden)
		case int32(pb.Code_Conflict):
			w.WriteHeader(http.StatusConflict)
		}

		w.Write(jsonResp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
