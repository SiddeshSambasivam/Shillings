package main

import (
	"net"

	"github.com/SiddeshSambasivam/shillings/pkg/models"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
)

func (env *DataEnv) createUser(conn net.Conn, req *pb.RequestSignup) {}

func (env *DataEnv) getUserProfile(conn net.Conn, req *pb.RequestGetUser) {

	row := env.DB.QueryRow("SELECT * FROM profiles WHERE user_id = ?", req.GetUserId())

	var u models.User
	err := row.Scan(
		&u.User_id,
		&u.First_name,
		&u.Middle_name,
		&u.Last_name,
		&u.Email,
		&u.Phone,
		&u.Balance,
		&u.Created_at,
		&u.Updated_at,
	)

	if err != nil {
		sendCmdErrResponse(
			conn,
			pb.Code_INTERNAL_SERVER_ERROR,
			"Error fetching user: "+err.Error(),
		)
	}

}
