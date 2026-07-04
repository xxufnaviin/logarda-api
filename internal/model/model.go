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

type MetricsLogsRequest struct {
	Username string `json:"username"`
	Duration string `json:"duration"`
}

type LogStatsRequest struct {
	Username      string `json:"username"`
	AggregateFunc string `json:"aggregateFunc"`
	AggregateCol  string `json:"aggregateCol"`
}

type ModelPerformanceRequest struct {
	Model string `json:"model"`
}
type PredictMetricsResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type LLMInferenceResponse struct {
	Data   LLMInfereceData `json:"data"`
	Status string          `json:"status"`
}

type LLMInfereceData struct {
	Explanation string `json:"explanation"`
	Solution    any    `json:"solution"`
}
type AWSErrorEvent struct {
	EventTime    time.Time `json:"eventTime"`
	ErrorCode    string    `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
	ServiceName  string    `json:"serviceName"`
	EventName    string    `json:"eventName"`
	Username     string    `json:"username"`
}

type User struct { // database schema
	Username        string
	Password        string
	AccessKeyID     *string
	AccessKeySecret *string
	Region          *string
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

type PredictedMetrics struct {
	MetricTime time.Time `json:"metrictime"`
	InstanceID string    `json:"instanceid"`
	Cpu        float64   `json:"cpu"`
	Network    float64   `json:"network"`
	Memory     float64   `json:"memory"`
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

type LogStats struct {
	Column   string `json:"column"`
	AggValue string `json:"aggvalue"`
}
type Message struct {
	MsgType string `json:"type"`
	Msg     any    `json:"msg"` // can be Logs or Metric
}

// map user id to channels for messages
var OnlineUsers = make(map[string]chan Message)
