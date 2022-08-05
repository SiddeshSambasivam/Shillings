package main

import (
	"net"

	protocols "github.com/SiddeshSambasivam/shillings/pkg/protocols"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"google.golang.org/protobuf/proto"
)

func readCommand(conn net.Conn, requestPb *pb.RequestCommand) error {

	dataBuffer, readErr := protocols.ReadProtocolData(conn)
	if readErr != nil {
		sendCmdErrResponse(
			conn,
			pb.Code_DATA_LOSS,
			"Error reading header: "+readErr.Error(),
		)
		return readErr
	}

	err := proto.Unmarshal(dataBuffer, requestPb)
	if err != nil {
		sendCmdErrResponse(
			conn,
			pb.Code_DATA_LOSS,
			"Error unmarshalling: "+err.Error(),
		)
		return err
	}

	return nil
}
