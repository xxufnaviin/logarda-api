package model

import "time"

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

type AWSErrorEvent struct {
	EventTime    time.Time `json:"eventTime"`
	ErrorCode    string    `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
	ServiceName  string    `json:"serviceName"`
	EventName    string    `json:"eventName"`
	Username     string    `json:"username"`
}

type WebsocketRequest struct {
	Username string `json:"username"`
}

type Message struct {
	MsgType string `json:"type"`
	Msg     string `json:"msg"`
}

// map user id to channels for messages
var OnlineUsers = make(map[string]chan Message)
