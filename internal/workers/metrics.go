package workers

import (
	"fmt"
	"log"
	"logarda/internal/db"
	"logarda/internal/model"
	"logarda/utils"
	"time"
)

// var metrics model.Metrics

// @ goroutine
func MetricStreamWorker() {
	for {
		metrics := model.Metrics{}
		// consume metrics from redis queue
		metric, err := db.ConsumeMetricEvents()
		fmt.Println(metric)
		if err != nil {
			if err == db.RedisNil {
				log.Printf("Metrics worker timeout. Reconnecting.")
				continue
			}
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
		time.Sleep(1000000000) // for simulation
	}

}
