package main

import "github.com/SiddeshSambasivam/shillings/pkg/models"

type Status struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type SignUpResponse struct {
	Status Status `json:"status"`
}

type Auth struct {
	Token          string `json:"token"`
	ExpirationTime int64  `json:"expiration_time"`
}

type LoginResponse struct {
	Status Status `json:"status"`
}

type UserResponse struct {
	User   models.User `json:"user"`
	Status Status      `json:"status"`
}

type PaymentResponse struct {
	Transaction_id int64  `json:"transaction_id"`
	Status         Status `json:"status"`
}
