package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	User_id     int32   `json:"user_id"`
	First_name  string  `json:"first_name"`
	Middle_name string  `json:"middle_name"`
	Last_name   string  `json:"last_name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Balance     float32 `json:"balance"`
	Created_at  int64   `json:"created_at"`
	Updated_at  int64   `json:"updated_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CredentialData struct {
	Credential_id int32  `json:"credential_id"`
	User_id       int32  `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Created_at    int64  `json:"created_at"`
	Updated_at    int64  `json:"updated_at"`
}

type Claims struct {
	User_id            int32 `json:"user_id"`
	jwt.StandardClaims       // jwt.StandardClaims is a struct that contains the standard claims used by JWT.
}
