package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/SiddeshSambasivam/shillings/pkg/db"
	envtools "github.com/SiddeshSambasivam/shillings/pkg/env"
	errors "github.com/SiddeshSambasivam/shillings/pkg/errors"
	redisCache "github.com/SiddeshSambasivam/shillings/pkg/redis"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"github.com/go-redis/redis/v8"
)

var jwtKey []byte

type DataEnv struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (env *DataEnv) handleConnection(conn net.Conn) {

	defer conn.Close()

	requestPb := &pb.RequestCommand{}
	err := readCommand(conn, requestPb)
	if err != nil {
		log.Println("Error reading command: ", err)
		sendCmdErrResponse(
			conn,
			pb.Code_BAD_REQUEST,
			"Error reading command: "+err.Error(),
		)
	}

	cmd := requestPb.GetCommand()

	switch cmd {
	case pb.Command_LGN:
		env.commandHandlerLGN(conn)

	case pb.Command_SGN:
		env.commandHandlerSGN(conn)

	case pb.Command_USR:
		env.commandHandlerUSR(conn)

	case pb.Command_PAY:
		env.commandHandlerPAY(conn)

	case pb.Command_TPU:
		env.commandHandlerTPU(conn)

	case pb.Command_TXQ:
		env.commandHandlerTXN(conn)

	default:
		sendCmdErrResponse(
			conn,
			pb.Code_BAD_REQUEST,
			"Invalid Command: "+err.Error(),
		)
	}
}

func main() {

	var PORT = ":8080"
	var ADDR = "127.0.0.1"

	envPort := ":" + envtools.GetEnvVar("APP_PORT")
	log.Println("Loaded env var(port): ", envPort)
	if envPort != "" {
		PORT = envPort
	}

	envJwt := envtools.GetEnvVar("JWT_KEY")
	log.Println("Loaded env var(Jwt key): ", jwtKey)
	if envJwt != "" {
		jwtKey = []byte(envJwt)
	}

	ADDR = ADDR + PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp", PORT)
	errors.HandleErrorWithExt(err)

	log.Println("Serving application server @ : " + ADDR)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	errors.HandleErrorWithExt(err)
	defer listener.Close()

	db := db.InitDB()
	redis := redisCache.InitRedis()
	env := &DataEnv{DB: db, Redis: redis}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}

		go env.handleConnection(conn)

	}

}
