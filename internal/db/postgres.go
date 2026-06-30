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

	if err != nil {
		if err == NoRows { // return username exist = false if no rows found, and error = nil
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

func GetMetrics(ctx context.Context, username string, duration string) ([]model.Metrics, error) {
	var metrics []model.Metrics

	// get all metrics for the given duration (all instance)
	results, err := DB.Query(ctx, "SELECT * FROM metrics WHERE username = $1 AND metricTime >= NOW() - make_interval(hours => $2) ORDER BY instanceID, metricTime",
		username, duration)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	// scan each row into a struct and append to array
	for results.Next() {
		var metric model.Metrics
		err := results.Scan(&metric.MetricTime, &metric.InstanceID, &metric.Cpu, &metric.Network, &metric.Memory, &metric.Username)
		if err != nil {
			return nil, err
		}
		// append each row into list of metrics
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func GetErrorLogs(ctx context.Context, username string, duration string) ([]model.Logs, error) {
	var logs []model.Logs

	// get all logs for the given duration
	results, err := DB.Query(ctx, "SELECT * FROM logs WHERE username = $1 AND errorExplained = true AND eventTime >= NOW() - make_interval(hours => $2) ORDER BY eventTime",
		username, duration)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	// scan each row into a struct and append to array
	for results.Next() {
		var log model.Logs
		err := results.Scan(&log.EventTime, &log.ErrorCode, &log.ErrorMessage, &log.ServiceName, &log.EventName, &log.Username, &log.Explanation, &log.ErrorExplained)
		if err != nil {
			return nil, err
		}
		// append each row into list of logs
		logs = append(logs, log)
	}

	return logs, nil
}

func GetUser(ctx context.Context, username string, user *model.User) error {
	err := DB.QueryRow(ctx, "SELECT * FROM stg_users WHERE username=$1;", username).Scan(
		&user.Username, &user.Password, &user.AccessKeyID, &user.AccessKeySecret, &user.Region, &user.CollectorOn)

	return err
}

func UpdataCollectorToggle(ctx context.Context, username string) error {
	// current toggle between prd and stg here
	_, err := DB.Exec(ctx, "UPDATE stg_users SET collector_on = True WHERE username = $1;", username)

	return err
}

func GetErrorStats(ctx context.Context, username, aggregateFunc, aggregateCol string) ([]model.LogStats, error) {
	var logstats []model.LogStats

	query := fmt.Sprintf("SELECT %s, %s(*) FROM logs WHERE username ='%s' GROUP BY %s",aggregateCol, aggregateFunc, username, aggregateCol)
	// get all logs for the given duration
	results, err := DB.Query(ctx, query)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	// scan each row into a struct and append to array
	for results.Next() {
		var stats model.LogStats
		err := results.Scan(&stats.Column, &stats.AggValue)
		if err != nil {
			return nil, err
		}
		// append each row into list of logs
		logstats = append(logstats, stats)
	}

	return logstats, nil

}
