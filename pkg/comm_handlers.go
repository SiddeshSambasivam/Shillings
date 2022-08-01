package pkg

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"google.golang.org/protobuf/proto"
)

func SendCmdRequest(client net.Conn, req *pb.RequestCommand) {

	data, err := proto.Marshal(req)
	if err != nil {
		log.Println("Error marshalling:", err)
		return
	}

	// First send the len of the body
	b := make([]byte, 4) // 4 bytes for header.
	binary.BigEndian.PutUint32(b, uint32(len(data)))

	_, err = client.Write(b)
	if err != nil {
		log.Println("Error writing header:", err)
		return
	}

	_, err = client.Write(data)
	if err != nil {
		log.Println("Error writing data:", err)
		return
	}

}

func SendCmdResponse(conn net.Conn, resp *pb.ResponseCommand) {

	data, err := proto.Marshal(resp)
	if err != nil {
		log.Println("Error marshalling response: ", err)
		return
	}

	b := make([]byte, 4) // 4 bytes for header.
	binary.BigEndian.PutUint32(b, uint32(len(data)))

	_, err = conn.Write(b)
	if err != nil {
		log.Println("Error writing header: ", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("Error writing data: ", err)
		return
	}

}
