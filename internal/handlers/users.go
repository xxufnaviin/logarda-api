package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"logarda/internal/db"
	"logarda/internal/model"
	"logarda/utils"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request model.LoginRequest
	var user model.User
	var hashedPassword string

	// decode the request body into user struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest) // return error if not able to parse body
		return
	}

	hashedPassword = utils.HashString(request.Password) // hash password before comparison

	err = db.VerifyUserCredentials(ctx, request.Username, hashedPassword, &user) //check if credentials exist in database

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if err == db.NoRows {
			fmt.Println("Invalid credentials!")
			json.NewEncoder(w).Encode(map[string]any{
				"message": "login failed",
				"status":  404,
				"error":   "Invalid credentials.",
			})
			return
		} else {
			json.NewEncoder(w).Encode(map[string]any{
				"message": "login failed",
				"status":  http.StatusBadGateway,
				"error":   err.Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "login success",
		"status":  http.StatusOK,
		"data": map[string]string{
			"username": user.Username,
		},
	})

	fmt.Println("Login success")
}

func Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request model.RegisterRequest
	var userExists bool

	// decode the request body into register struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest) // return error if not able to parse body
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// check if username exists already in database
	err, userExists = db.CheckUniqueUsername(ctx, request.Username)

	if err != nil {
		fmt.Println("Failed to check if username exists. Register Failed.")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "register failed",
			"status":  404,
			"error":   err.Error(),
		})
		return
	}

	if userExists {
		fmt.Println("Username already exists! Please try again.")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "register failed",
			"status":  404,
			"error":   "username exists",
		})
		return
	}

	// hash password before storing
	hashedPassword := utils.HashString(request.Password)

	// register new user and update database
	err = db.RegisterNewUser(ctx, request.Username, hashedPassword)

	if err != nil {
		fmt.Println("Invalid username or password!")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "register failed",
			"status":  404,
			"error":   "Invalid credentials.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "register success",
		"status":  http.StatusOK,
		"data": map[string]string{
			"username": request.Username,
		},
	})

	fmt.Println("User Registered")

}

func SaveAWSCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var request model.AWSCredentialsRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest) // return error if not able to parse body
		return
	}
	// check if keys are in valid format
	if !utils.IsValidAccessKey(request.AccessKeyID) {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "failed",
			"status":  http.StatusNotFound,
			"error":   "Access Key ID is not in valid format!",
		})
		return
	}
	if !utils.IsValidSecretKey(request.AccessKeySecret) {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "failed",
			"status":  http.StatusNotFound,
			"error":   "Access Key Secret is not in valid format!",
		})
		return
	}

	err = utils.VerifyAWSCredentials(request.AccessKeyID, request.AccessKeySecret, request.Region) // verify with AWS if user exist
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "failed",
			"status":  http.StatusNotFound,
			"error":   "AWS keys provided are invalid!",
		})
		return
	}

	hasAllPermissions := utils.HasAllPermissions(request.AccessKeyID, request.AccessKeySecret, request.Region)
	if !hasAllPermissions{
		json.NewEncoder(w).Encode(map[string]any{
			"message": "failed",
			"status":  http.StatusNotFound,
			"error":   "AWS access keys have insufficient permissions! Please enable EC2, CloudWatch and CloudTrail",
		})
		return
	}

	// encrypt keys before saving to database
	encryptedID := utils.EncryptString(request.AccessKeyID)
	encryptedSecret := utils.EncryptString(request.AccessKeySecret)

	err = db.SaveAWSCredentials(ctx, encryptedID, encryptedSecret, request.Region, request.Username) //save aws credentials in database

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "failed",
			"status":  http.StatusNotFound,
			"error":   err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "update success",
		"status":  http.StatusOK,
	})

}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	var user model.User
	username := r.URL.Query().Get("username")
	ctx := r.Context()

	err := db.GetUser(ctx, username, &user) // get user from database

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Failed. User not found",
			"status":  http.StatusNotFound,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"data":   user,
		"status": http.StatusOK,
	})

}

func SetCollectorOn(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	ctx := r.Context()

	err := db.UpdataCollectorToggle(ctx, username) // set collector toggle on

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Failed. Collector not turned on.",
			"status":  http.StatusNotFound,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"message":   "Collector turned on.",
		"status": http.StatusOK,
	})

}
