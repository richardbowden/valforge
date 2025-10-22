package main

type User struct {
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"gte=18"`
	Pwd1  string `json:"pwd1" validate:"minlen=6"`
	Pwd2  string `json:"pwd2" validate:"eqfieldsecure=Pwd1"`
	Email string `json:"email" validate:"email"`
	Color string `json:"color" validate:"required"`
}
