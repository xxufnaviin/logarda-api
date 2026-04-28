package handlers

import (
	"net/http"	
	"logarda/internal/config"
)

func Login(w http.ResponseWriter, r http.Request){
	// check if database got user
}


func SetAWSCredentials(w http.ResponseWriter, r http.Request){
	config.AWSAccessKeyID = "aws"
	config.AWSAccessKeySecret = "aws"
}
