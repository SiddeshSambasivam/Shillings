package protocols

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ProcessDataBuffer reads the data buffer from the connection.
// It returns the data buffer and an error if any.
//
// Specifications:
// - Header is four bytes long and is the length of the payload.
// - Memory is allocated for the data buffer from the header
func ReadProtocolData(conn net.Conn) ([]byte, error) {
	var headerByteSize int = 4
	headerBuffer := make([]byte, headerByteSize)

	_, readErr := conn.Read(headerBuffer)
	if readErr != nil && readErr != io.EOF {
		return nil, readErr
	}

	dataByteSize := int(binary.BigEndian.Uint32(headerBuffer))
	dataBuffer := make([]byte, dataByteSize)
	_, dataReadError := conn.Read(dataBuffer)

	if dataReadError != nil &&
		dataReadError != io.EOF &&
		len(dataBuffer) != dataByteSize {
		return nil, dataReadError
	}

	return dataBuffer, nil
}

func SendProtocolData(client net.Conn, req protoreflect.ProtoMessage) error {

	data, err := proto.Marshal(req)
	if err != nil {
		log.Println("Error marshalling:", err)
		return err
	}

	// First send the len of the body
	b := make([]byte, 4) // 4 bytes for header.
	binary.BigEndian.PutUint32(b, uint32(len(data)))

	_, err = client.Write(b)
	if err != nil {
		log.Println("Error writing header:", err)
		return err
	}

	_, err = client.Write(data)
	if err != nil {
		log.Println("Error writing data:", err)
		return err
	}

	return nil

}
