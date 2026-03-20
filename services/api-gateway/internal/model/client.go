package model

import "errors"

var ErrNotFound = errors.New("not found")

type Client struct {
	ID      int32  `json:"id"`
	Surname string `json:"surname"`
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Email   string `json:"email"`
}

type CreateClientRequest struct {
	Surname string `json:"surname" binding:"required"`
	Name    string `json:"name"    binding:"required"`
	Age     int32  `json:"age"     binding:"required,min=1,max=150"`
	Email   string `json:"email"   binding:"required,email"`
}
