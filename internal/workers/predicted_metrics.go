package workers

import (
	"fmt"
	"logarda/internal/db"
	"logarda/internal/model"
	"logarda/utils"
	"time"
)

var predicted_metric string
var predicted_metrics model.PredictedMetrics

// @ goroutine
func PredictedMetricStreamWorker() {
	for {
		// consume metrics from redis queue
		predicted_metric, err = db.ConsumePredictedMetricEvents()
		fmt.Println(predicted_metric)
		if err != nil {
			fmt.Printf("Error getting event data.")
			continue
		}

		// unmarshal string to JSON
		err = utils.UnmarshalAWSPredictedMetricEvent(predicted_metric, &predicted_metrics)
		if err != nil {
			fmt.Printf("Error during parsing event data.")
			continue
		}

		// stream to websocket (online users)
		mu.Lock() // get the msg channel mutually exclusive to prevent unsafe actions
		msgChannel, ok := model.OnlineUsers[predicted_metrics.Username]
		mu.Unlock()

		if ok {
			msgChannel <- model.Message{
				MsgType: "predicted_metrics",
				Msg:     predicted_metrics}
		}
		time.Sleep(500000000) // for simulation
	}

}
