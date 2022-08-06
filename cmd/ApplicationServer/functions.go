package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/SiddeshSambasivam/shillings/pkg/models"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (env *DataEnv) createUser(req *pb.RequestSignup) error {

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

	createCredQry := "INSERT INTO credentials (user_id, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	ins, err = env.DB.PrepareContext(ctx, createCredQry)
	if err != nil {
		log.Println("Error preparing create cred query", err)
	}

	res, err = ins.ExecContext(ctx, user_id, email, password, created_at, updated_at)
	if err != nil {
		log.Println("Error inserting cred", err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Println("Error getting last inserted id", err)
	}

	return nil
}

func (env *DataEnv) loginUser(req *pb.RequestLogin) (protoreflect.ProtoMessage, error) {

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

	err = bcrypt.CompareHashAndPassword(
		[]byte(creds.Password),
		[]byte(req.Credentials.GetPassword()),
	)
	if err != nil {
		err = errors.New("incorrect password")
		return nil, err
	}

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

func (env *DataEnv) GetUser(req *pb.RequestGetUser) (protoreflect.ProtoMessage, error) {

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

func (env *DataEnv) PayUser(req *pb.RequestPayUser) (protoreflect.ProtoMessage, error) {

	isauth, claims, err := isAuthenticated(req.GetAuth().GetToken())
	if err != nil {
		return nil, err
	}
	if !isauth {
		return nil, errors.New("user not authenticated")
	}

	exists, err := env.checkUserExists(req.GetReceiverEmail())
	if err != nil {
		log.Println("Error checking if user exists", err)
		return nil, err
	}

	if !exists {
		err = errors.New("Sender with email " + req.GetReceiverEmail() + " does not exist")
		return nil, err
	}

	sender_id := claims.User_id
	sender, err := env.getUserAccountByID(sender_id)
	if err != nil {
		log.Println("Error fetching user", err)
		return nil, err
	}

	receiver, err := env.getUserAccountByEmail(req.GetReceiverEmail())
	if err != nil {
		log.Println("Error fetching user", err)
		return nil, err
	}

	if receiver.User_id == sender_id {
		err = errors.New("sender and receiver cannot be the same")
		return nil, err
	}

	if sender.Balance < req.GetAmount() {
		err = errors.New("insufficient balance")
		return nil, err
	}

	sender.Balance -= req.GetAmount()
	receiver.Balance += req.GetAmount()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	tx, err := env.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error beginning transaction", err)
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = ? WHERE user_id = ?", sender.Balance, sender.User_id)
	if err != nil {
		log.Println("Error updating sender balance", err)
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = ? WHERE user_id = ?", receiver.Balance, receiver.User_id)
	if err != nil {
		log.Println("Error updating receiver balance", err)
		return nil, err
	}

	// insert new transaction to transactions table
	res, err := tx.ExecContext(ctx,
		"INSERT INTO transactions (sender_id, receiver_id, amount, created_at, sender_email, receiver_email) VALUES (?, ?, ?, ?, ?, ?)",
		sender.User_id, receiver.User_id, req.GetAmount(), time.Now().Unix(), sender.Email, receiver.Email)
	if err != nil {
		log.Println("Error inserting transaction", err)
		return nil, err
	}

	transaction_id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert id", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction", err)
		return nil, err
	}

	resp := &pb.ResponsePayUser{
		Status: &pb.Status{
			Code:    pb.Code_OK,
			Message: "Payment successful",
		},
		TransactionId: transaction_id,
	}

	return resp, nil

}

func (env *DataEnv) TopUpUser(req *pb.RequestTopupUser) (protoreflect.ProtoMessage, error) {

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

	// begin a db transaction
	tx, err := env.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error beginning transaction", err)
		return nil, err
	}

	// update the user's balance
	_, err = tx.ExecContext(
		ctx,
		"UPDATE users SET balance = ? WHERE user_id = ?",
		req.GetAmount()+user.Balance,
		claims.User_id,
	)
	if err != nil {
		log.Println("Error updating user balance", err)
		return nil, err
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO transactions (sender_id, receiver_id, amount, created_at, sender_email, receiver_email) VALUES (?, ?, ?, ?, ?, ?)",
		claims.User_id, claims.User_id, req.GetAmount(), time.Now().Unix(), user.Email, user.Email)
	if err != nil {
		log.Println("Error inserting transaction", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction", err)
		return nil, err
	}

	resp := &pb.ResponseTopupUser{
		Status: &pb.Status{
			Code:    pb.Code_OK,
			Message: "Topup successful",
		},
	}

	return resp, nil

}

func (env *DataEnv) GetUSRTransactions(req *pb.RequestGetUserTransactions) (protoreflect.ProtoMessage, error) {

	isauth, claims, err := isAuthenticated(req.GetAuth().GetToken())
	if err != nil {
		return nil, err
	}
	if !isauth {
		return nil, errors.New("user not authenticated")
	}

	log.Println(claims.Id)
	// get all transactions for the user
	var transactions []*pb.Transaction
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	rows, err := env.DB.QueryContext(
		ctx,
		"SELECT * from transactions where sender_id = ? or receiver_id = ?",
		claims.User_id, claims.User_id,
	)
	if err != nil {
		log.Println("Error fetching transactions", err)
		return nil, err
	}

	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.Transaction_id, &transaction.Sender_id, &transaction.Receiver_id, &transaction.Amount, &transaction.Created_at, &transaction.Sender_email, &transaction.Receiver_email)
		if err != nil {
			log.Println("Error scanning transaction", err)
		}
		transactions = append(transactions, &pb.Transaction{
			TransactionId: transaction.Transaction_id,
			SenderId:      transaction.Sender_id,
			ReceiverId:    transaction.Receiver_id,
			Amount:        transaction.Amount,
			CreatedAt:     transaction.Created_at,
			SenderEmail:   transaction.Sender_email,
			ReceiverEmail: transaction.Receiver_email,
		})
	}

	resp := &pb.ResponseGetUserTransactions{
		Status: &pb.Status{
			Code:    pb.Code_OK,
			Message: "Transactions fetched successfully",
		},
		Transactions: transactions,
	}

	return resp, nil
}
