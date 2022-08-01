package main

import (
	"net"

	pkg "github.com/SiddeshSambasivam/shillings/pkg"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
)

func sendCmdErrResponse(conn net.Conn, status_code pb.Code, err_message string) {
	response := &pb.ResponseCommand{
		Status: &pb.Status{
			Code:    status_code,
			Message: err_message,
		},
	}

	pkg.SendCmdResponse(conn, response)
}
