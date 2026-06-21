package workers

import (
	"fmt"
	"logarda/internal/db"
	"logarda/internal/model"
	"logarda/utils"
	"time"
)

var metric string
var metrics model.Metrics

// @ goroutine
func MetricStreamWorker() {
	for {
		// consume metrics from redis queue
		metric, err = db.ConsumeMetricEvents()
		fmt.Println(metric)
		if err != nil {
			fmt.Printf("Error getting event data.")
			continue
		}

		// unmarshal string to JSON
		err = utils.UnmarshalAWSMetricEvent(metric, &metrics)
		if err != nil {
			fmt.Printf("Error during parsing event data.")
			continue
		}

		// stream to websocket (online users)
		mu.Lock() // get the msg channel mutually exclusive to prevent unsafe actions
		msgChannel, ok := model.OnlineUsers[metrics.Username]
		mu.Unlock()

		if ok {
			msgChannel <- model.Message{
				MsgType: "metrics",
				Msg:     metrics}
		}
		time.Sleep(500000000) // for simulation
	}

}
