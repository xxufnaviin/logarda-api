package model

// define structs here

type User struct {
	Username string
	Password string
}
type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
