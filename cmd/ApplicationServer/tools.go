package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/SiddeshSambasivam/shillings/pkg/models"
	protocols "github.com/SiddeshSambasivam/shillings/pkg/protocols"
	"github.com/SiddeshSambasivam/shillings/proto/shillings/pb"
	"github.com/dgrijalva/jwt-go"
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

func sendCmdErrResponse(conn net.Conn, status_code pb.Code, err_message string) {
	response := &pb.ResponseCommand{
		Status: &pb.Status{
			Code:    status_code,
			Message: err_message,
		},
	}

	protocols.SendProtocolData(conn, response)
}

func SendSignupErrResponse(conn net.Conn, status_code pb.Code, err_message string) {
	response := &pb.ResponseSignup{
		Status: &pb.Status{
			Code:    status_code,
			Message: err_message,
		},
	}

	protocols.SendProtocolData(conn, response)
}

func SendLoginErrResponse(conn net.Conn, status_code pb.Code, err_message string) {
	response := &pb.ResponseLogin{
		Status: &pb.Status{
			Code:    status_code,
			Message: err_message,
		},
	}

	protocols.SendProtocolData(conn, response)
}

func SendUserErrResponse(conn net.Conn, status_code pb.Code, err_message string) {
	response := &pb.ResponseGetUser{
		Status: &pb.Status{
			Code:    status_code,
			Message: err_message,
		},
	}

	protocols.SendProtocolData(conn, response)
}

func (env *DataEnv) checkUserExists(email string) (bool, error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	var exists bool

	row := env.DB.QueryRowContext(ctx, "SELECT EXISTS (SELECT EMAIL FROM users WHERE EMAIL = ?)", email)
	if err := row.Scan(&exists); err != nil {
		log.Println("Error checking if user exists", err)
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}

func (env *DataEnv) getUserAccountByID(user_id int32) (models.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	var user models.User
	fetchCredQuery := "SELECT * from users where user_id = ?"
	err := env.DB.QueryRowContext(
		ctx,
		fetchCredQuery, user_id,
	).Scan(&user.User_id, &user.First_name, &user.Middle_name, &user.Last_name, &user.Email, &user.Phone, &user.Balance, &user.Created_at, &user.Updated_at)

	if err != nil {
		log.Println("Error fetching user", err)
		return user, err
	}

	return user, nil

}

func (env *DataEnv) getUserAccountByEmail(email string) (models.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	var user models.User
	fetchCredQuery := "SELECT * from users where email = ?"
	err := env.DB.QueryRowContext(
		ctx,
		fetchCredQuery, email,
	).Scan(&user.User_id, &user.First_name, &user.Middle_name, &user.Last_name, &user.Email, &user.Phone, &user.Balance, &user.Created_at, &user.Updated_at)

	if err != nil {
		log.Println("Error fetching user", err)
		return user, err
	}

	return user, nil

}

func generateJWT(user_id int32) (string, time.Time, error) {

	// Generate JWT token
	expirationTime := time.Now().Add(15 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		User_id: user_id,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, expirationTime, err
}

func isAuthenticated(token string) (bool, models.Claims, error) {

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, *claims, err
	}

	if !tkn.Valid {
		return false, *claims, err
	}

	return true, *claims, err
}
