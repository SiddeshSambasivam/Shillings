package main

type Status struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type SignUpResponse struct {
	Status Status `json:"status"`
}
