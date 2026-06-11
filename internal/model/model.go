package model

// define structs here

type User struct { // database schema
	Username        string
	Password        string
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	CollectorOn     bool
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AWSCredentialsRequest struct {
	Username        string `json:"username"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	Region          string `json:"region"`
}
