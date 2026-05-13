package handlers

import (
	"encoding/json"
	"fmt"
	"logarda/internal/db"
	"logarda/internal/model"
	"logarda/utils"
	"net/http"

	// "context"
)

func Login(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	var request model.RequestLogin
	var user model.User
	var hashedPassword string

	// decode the request body into user struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest) // return error if not able to parse body
		return
	}

	hashedPassword = utils.HashString(request.Password) // hash password before comparison

	err = db.DB.QueryRow(ctx,"SELECT username, password FROM users WHERE username=$1 AND password=$2;", 
						request.Username, hashedPassword).Scan(&user.Username, &user.Password)					

	if err != nil{
		if err == db.NoRows{
			fmt.Println("Invalid credentials!")
			json.NewEncoder(w).Encode(map[string]any{
				"message": "failed",
				"status": http.NotFound,
				"error": "Invalid credentials.",
			})
			return
		}
		fmt.Printf("%s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "login success",
		"status": http.StatusOK,
	})

	fmt.Println("Login success")	
}


