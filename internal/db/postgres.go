package db

import (
	"context"
	"fmt"
	"logarda/internal/config"
	"logarda/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var err error
var NoRows error

func CreatePostgreSQLPool() {
	ctx := context.Background()
	NoRows = pgx.ErrNoRows

	DB, err = pgxpool.New(ctx, config.POSTGRES_URL) // assigns pool to global variable
	if err != nil {
		DB = nil
		fmt.Println("Error connecting to Postgres database!")
		return
	}
}

func RegisterNewUser(ctx context.Context, username, hashedpassword string) error {
	// current toggle between prd and stg here
	_, err := DB.Exec(ctx, "INSERT INTO stg_users (username, password) VALUES ($1, $2);", username, hashedpassword)

	return err
}

func VerifyUserCredentials(ctx context.Context, username, hashedPassword string, user *model.User) error {
	// current toggle between prd and stg here
	err := DB.QueryRow(ctx, "SELECT username, password FROM stg_users WHERE username=$1 AND password=$2;",
		username, hashedPassword).Scan(&user.Username, &user.Password)
	return err
}

func SaveAWSCredentials(ctx context.Context, encryptedID, encryptedSecret, region, username string) error {
	// current toggle between prd and stg here
	_, err := DB.Exec(ctx, "UPDATE stg_users SET awskeyid = $1, awskeysecret=$2, awsregion=$3 WHERE username=$4;",
		encryptedID, encryptedSecret, region, username)
	return err
}

func CheckUniqueUsername(ctx context.Context, username string) (error, bool) {
	var existingUser string
	err := DB.QueryRow(ctx, "SELECT username FROM stg_users WHERE username=$1;",
		username).Scan(&existingUser)

	if err != nil{
		if err == NoRows{ // return username exist = false if no rows found, and error = nil
			return nil, false
		}
		return err, false // return error if found and handle
	}

	return err, true // return username exist = true if rows found
}

func SaveErrorExplanations(ctx context.Context, event *model.AWSErrorEvent, explanation string) error {
	// current toggle between prd and stg here
	_, err := DB.Exec(ctx, "UPDATE stg_logs SET explanation = $1, errorExplained=$2  WHERE eventTime = $3 AND errorCode = $4 AND errorMessage = $5;",
		explanation, true, event.EventTime, event.ErrorCode, event.ErrorMessage)
	return err
}