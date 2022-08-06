package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"golang.org/x/crypto/bcrypt"
)

func (env *DataEnv) createUser(conn net.Conn, req *pb.RequestSignup) error {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	var exists bool

	row := env.DB.QueryRowContext(ctx, "SELECT EXISTS (SELECT EMAIL FROM users WHERE EMAIL = ?)", req.GetUser().GetEmail())
	if err := row.Scan(&exists); err != nil {
		log.Println("Error checking if user exists", err)
		return err
	} else if exists {
		err = errors.New("User with email " + req.GetUser().GetEmail() + " already exists")
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetCredentials().GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password", err)
		return err
	}

	firstName := req.GetUser().GetFirstName()
	middleName := req.GetUser().GetMiddleName()
	lastName := req.GetUser().GetLastName()
	email := req.GetUser().GetEmail()
	phone := req.GetUser().GetPhone()
	balance := req.GetUser().GetBalance()
	password := string(hash)
	created_at := time.Now().Unix()
	updated_at := time.Now().Unix()

	createUsrQry := "INSERT INTO users (first_name, middle_name, last_name, email, phone, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	ins, err := env.DB.PrepareContext(ctx, createUsrQry)
	if err != nil {
		log.Println("Error preparing create user query", err)
	}
	res, err := ins.ExecContext(ctx, firstName, middleName, lastName, email, phone, balance, created_at, updated_at)
	if err != nil {
		log.Println("Error inserting user", err)
	}

	user_id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error getting last inserted id", err)
	}

	log.Println("Created User_id:", user_id)

	createCredQry := "INSERT INTO credentials (user_id, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	ins, err = env.DB.PrepareContext(ctx, createCredQry)
	if err != nil {
		log.Println("Error preparing create cred query", err)
	}

	res, err = ins.ExecContext(ctx, user_id, email, password, created_at, updated_at)
	if err != nil {
		log.Println("Error inserting cred", err)
	}

	cred_id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error getting last inserted id", err)
	}

	log.Println("Created Credential_id:", cred_id)

	return nil
}
