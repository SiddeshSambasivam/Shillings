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

	err = env.createUser(conn, &req)
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
