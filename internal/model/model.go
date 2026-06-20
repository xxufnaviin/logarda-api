package model

import "time"

// define structs here
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
type WebsocketRequest struct {
	Username string `json:"username"`
}

type AWSErrorEvent struct {
	EventTime    time.Time `json:"eventTime"`
	ErrorCode    string    `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
	ServiceName  string    `json:"serviceName"`
	EventName    string    `json:"eventName"`
	Username     string    `json:"username"`
}

type MetricsLogsRequest struct {
	Username string `json:"username"`
	Duration string `json:"duration"`
}

type User struct { // database schema
	Username        string
	Password        string
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	CollectorOn     bool
}

type Metrics struct {
	MetricTime time.Time `json:"metricTime"`
	InstanceID string    `json:"instanceID"`
	Cpu        float32   `json:"cpu"`
	Network    float32   `json:"network"`
	Memory     float32   `json:"memory"`
	Username   string    `json:"username"`
}

type Logs struct {
	EventTime      time.Time `json:"eventTime"`
	ErrorCode      string    `json:"errorCode"`
	ErrorMessage   string    `json:"errorMessage"`
	ServiceName    string    `json:"serviceName"`
	EventName      string    `json:"eventName"`
	Username       string    `json:"username"`
	Explanation    string    `json:"explanation"` 
	ErrorExplained bool      `json:"errorExplained"`
}
type Message struct {
	MsgType string `json:"type"`
	Msg     any    `json:"msg"` // can be Logs or Metric
}

// map user id to channels for messages
var OnlineUsers = make(map[string]chan Message)
