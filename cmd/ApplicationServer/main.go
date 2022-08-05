package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"

	pkg "github.com/SiddeshSambasivam/shillings/pkg"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"google.golang.org/protobuf/proto"
)

var PORT = ":8080"
var ADDR = "127.0.0.1"

func readCommand(conn net.Conn, requestPb *pb.RequestCommand) error {

	var headerByteSize int = 4
	headerBuffer := make([]byte, headerByteSize)

	_, readErr := conn.Read(headerBuffer)

	if readErr != nil && readErr != io.EOF {
		sendCmdErrResponse(
			conn,
			pb.Code_DATA_LOSS,
			"Error reading header: "+readErr.Error(),
		)

		return readErr
	}

	dataByteSize := int(binary.BigEndian.Uint32(headerBuffer))

	dataBuffer := make([]byte, dataByteSize)
	_, dataReadError := conn.Read(dataBuffer)

	if dataReadError != nil && dataReadError != io.EOF {
		sendCmdErrResponse(
			conn,
			pb.Code_DATA_LOSS,
			"Error reading data: "+dataReadError.Error(),
		)
		return dataReadError
	}

	if len(dataBuffer) != dataByteSize {
		sendCmdErrResponse(
			conn,
			pb.Code_DATA_LOSS,
			"Error missing data: "+dataReadError.Error(),
		)
		return dataReadError
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

func handleConnection(conn net.Conn) {

	defer conn.Close()

	requestPb := &pb.RequestCommand{}
	err := readCommand(conn, requestPb)
	if err != nil {
		log.Println("Error reading command: ", err)
		return
	}

	cmd := requestPb.GetCommand()

	switch cmd {
	case pb.Command_LGN:
		log.Println("Login command received")
	case pb.Command_SGN:
		log.Println("Signup command received")
	case pb.Command_USR:
		log.Println("User command received")
	case pb.Command_PAY:
		log.Println("Pay command received")
	case pb.Command_TPU:
		log.Println("Topup command received")
	case pb.Command_TXQ:
		log.Println("Transaction query command received")
	default:
		log.Println("Unknown command received")
	}
	// If the command is not supported, return an error
	// Create request and response for the command

	// 1. Read the payload
	// 2. redirect to the specific handler

}

func main() {

	envPort := ":" + pkg.GetEnvVar("APP_PORT")
	log.Println("Loaded env var: ", envPort)

	if envPort != "" {
		PORT = envPort
	}

	ADDR = ADDR + PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp", PORT)
	pkg.HandleErrorWithExt(err)

	log.Println("Serving application server @ : " + ADDR)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	pkg.HandleErrorWithExt(err)

	Db := pkg.DbConn()
	Db.SetConnMaxLifetime(3 * time.Hour)
	Db.SetMaxOpenConns(3000)
	Db.SetMaxIdleConns(3000)
	defer listener.Close()

	row, err := Db.Query("SELECT * FROM users")
	pkg.HandleErrorWithExt(err)
	// iterate over the rows
	for row.Next() {
		log.Println("Row: ", row)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}

		log.Println("Accepted connection: ", conn.RemoteAddr())
		go handleConnection(conn)

	}

}
