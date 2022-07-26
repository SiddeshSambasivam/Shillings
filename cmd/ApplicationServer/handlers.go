package main

import (
	"log"
	"net"

	protocols "github.com/SiddeshSambasivam/shillings/pkg/protocols"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"google.golang.org/protobuf/proto"
)

// Handles user signup command
func (env *DataEnv) commandHandlerSGN(conn net.Conn) {

	dataBuffer, err := protocols.ReadProtocolData(conn)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendSignupErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	req := pb.RequestSignup{}
	err = proto.Unmarshal(dataBuffer, &req)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendSignupErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	err = env.createUser(&req)
	if err != nil {
		log.Println("Error creating user: ", err)
		SendSignupErrResponse(
			conn,
			pb.Code_Conflict,
			"Error creating user: "+err.Error(),
		)
	}

	resp := &pb.ResponseSignup{
		Status: &pb.Status{
			Code:    pb.Code_OK,
			Message: "User created successfully",
		},
	}

	protocols.SendProtocolData(conn, resp)
}

func (env *DataEnv) commandHandlerLGN(conn net.Conn) {
	dataBuffer, err := protocols.ReadProtocolData(conn)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendLoginErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	req := pb.RequestLogin{}
	err = proto.Unmarshal(dataBuffer, &req)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendLoginErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	resp, err := env.loginUser(&req)
	if err != nil {
		log.Println("Error authenticating user: ", err)
		SendLoginErrResponse(
			conn,
			pb.Code_Conflict,
			"Error authenticating user: "+err.Error(),
		)
	}

	protocols.SendProtocolData(conn, resp)
}

func (env *DataEnv) commandHandlerUSR(conn net.Conn) {
	dataBuffer, err := protocols.ReadProtocolData(conn)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	req := pb.RequestGetUser{}
	err = proto.Unmarshal(dataBuffer, &req)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	resp, err := env.GetUser(&req)
	if err != nil {
		log.Println("Error fetching user: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_Conflict,
			"Error fetching user: "+err.Error(),
		)
	}

	protocols.SendProtocolData(conn, resp)
}

func (env *DataEnv) commandHandlerPAY(conn net.Conn) {
	dataBuffer, err := protocols.ReadProtocolData(conn)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	req := pb.RequestPayUser{}
	err = proto.Unmarshal(dataBuffer, &req)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	resp, err := env.PayUser(&req)
	if err != nil {
		log.Println("Error paying user: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_BAD_REQUEST,
			"Error paying user: "+err.Error(),
		)
	}

	protocols.SendProtocolData(conn, resp)
}

func (env *DataEnv) commandHandlerTPU(conn net.Conn) {
	dataBuffer, err := protocols.ReadProtocolData(conn)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	req := pb.RequestTopupUser{}
	err = proto.Unmarshal(dataBuffer, &req)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	resp, err := env.TopUpUser(&req)
	if err != nil {
		log.Println("Error topping up user: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_BAD_REQUEST,
			"Error topping up user: "+err.Error(),
		)
	}

	protocols.SendProtocolData(conn, resp)
}

func (env *DataEnv) commandHandlerTXN(conn net.Conn) {
	dataBuffer, err := protocols.ReadProtocolData(conn)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	req := pb.RequestGetUserTransactions{}
	err = proto.Unmarshal(dataBuffer, &req)
	if err != nil {
		log.Println("Error reading data: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error reading data: "+err.Error(),
		)
	}

	resp, err := env.GetUSRTransactions(&req)
	if err != nil {
		log.Println("Error getting user transactions: ", err)
		SendUserErrResponse(
			conn,
			pb.Code_BAD_REQUEST,
			"Error getting user transactions: "+err.Error(),
		)
	}

	protocols.SendProtocolData(conn, resp)
}
