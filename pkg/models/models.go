package models

type User struct {
	User_id     int32   `json:"user_id"`
	First_name  string  `json:"first_name"`
	Middle_name string  `json:"middle_name"`
	Last_name   string  `json:"last_name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Balance     float64 `json:"balance"`
	Created_at  string  `json:"created_at"`
	Updated_at  string  `json:"updated_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
