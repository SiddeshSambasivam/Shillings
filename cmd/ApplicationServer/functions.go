package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/SiddeshSambasivam/shillings/pkg/models"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (env *DataEnv) createUser(conn net.Conn, req *pb.RequestSignup) error {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	exists, err := env.checkUserExists(req.GetUser().GetEmail())
	if err != nil {
		log.Println("Error checking if user exists", err)
		return err
	}

	if exists {
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

func (env *DataEnv) loginUser(conn net.Conn, req *pb.RequestLogin) (protoreflect.ProtoMessage, error) {

	exists, err := env.checkUserExists(req.GetCredentials().GetEmail())
	if err != nil {
		log.Println("Error checking if user exists", err)
		return nil, err
	}

	if !exists {
		err = errors.New("User with email " + req.GetCredentials().GetEmail() + " does not exist")
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	var creds models.CredentialData
	fetchCredQuery := "SELECT user_id, email, password from credentials where email = ?"
	err = env.DB.QueryRowContext(
		ctx,
		fetchCredQuery, req.Credentials.GetEmail(),
	).Scan(&creds.User_id, &creds.Email, &creds.Password)
	if err != nil {
		log.Println("Error fetching creds", err)
		return nil, err
	}

	log.Println("Fetched creds:", creds.Email, creds.Password)
	err = bcrypt.CompareHashAndPassword(
		[]byte(creds.Password),
		[]byte(req.Credentials.GetPassword()),
	)
	if err != nil {
		err = errors.New("incorrect password")
		return nil, err
	}

	log.Println("Logged in user:", creds.Email)

	tokenString, expirationTime, err := generateJWT(creds.User_id)
	if err != nil {
		log.Println("Error generating JWT", err)
		return nil, err
	}
	resp := &pb.ResponseLogin{
		Auth: &pb.Auth{
			Token:          tokenString,
			ExpirationTime: expirationTime.Unix(),
		},
		Status: &pb.Status{
			Code:    pb.Code_OK,
			Message: "Login successful",
		},
	}

	return resp, nil
}

func (env *DataEnv) GetUser(conn net.Conn, req *pb.RequestGetUser) (protoreflect.ProtoMessage, error) {

	isauth, claims, err := isAuthenticated(req.GetAuth().GetToken())
	if err != nil {
		return nil, err
	}
	if !isauth {
		return nil, errors.New("user not authenticated")
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	var user models.User
	fetchCredQuery := "SELECT * from users where user_id = ?"
	err = env.DB.QueryRowContext(
		ctx,
		fetchCredQuery, claims.User_id,
	).Scan(&user.User_id, &user.First_name, &user.Middle_name, &user.Last_name, &user.Email, &user.Phone, &user.Balance, &user.Created_at, &user.Updated_at)
	if err != nil {
		log.Println("Error fetching user", err)
		return nil, err
	}

	log.Println("Fetched user: ", user.Email)

	u := &pb.User{
		UserId:     user.User_id,
		FirstName:  user.First_name,
		MiddleName: user.Middle_name,
		LastName:   user.Last_name,
		Email:      user.Email,
		Phone:      user.Phone,
		Balance:    user.Balance,
		CreatedAt:  user.Created_at,
		UpdatedAt:  user.Updated_at,
	}

	// TODO: add auth token to response
	resp := &pb.ResponseGetUser{
		User: u,
		Status: &pb.Status{
			Code:    pb.Code_OK,
			Message: "User fetched successfully",
		},
	}

	return resp, nil

}
