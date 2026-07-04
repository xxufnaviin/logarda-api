package workers

import (
	"context"
	"fmt"
	"log"
	"logarda/internal/db"
	"logarda/internal/handlers"
	"logarda/internal/model"
	"logarda/utils"
	"sync"
	"time"
)

var errorExplanation string

var errorEvent model.AWSErrorEvent
var mu sync.RWMutex

// @ goroutine
func ErrorLogsWorker() {
	for {
		ctx := context.Background()

		// consume error event from redis queue
		errorMsg, err := db.ConsumeErrorEvents()
		fmt.Println(errorMsg)
		if err != nil {
			if err == db.RedisNil {
				log.Printf("Logs worker timeout. Reconnecting.")
				continue
			}
			fmt.Printf("Error getting event data.")
			continue
		}

		// unmarshal string to JSON before saving to database
		err = utils.UnmarshalAWSErrorEvent(errorMsg, &errorEvent)
		if err != nil {
			fmt.Printf("Error during parsing event data.")
			continue
		}


		// make api call
		errorExplanation = handlers.GetLLMInference(&errorEvent) // placeholder

		// stream to websocket (online users)
		mu.Lock() // get the msg channel mutually exclusive to prevent unsafe actions
		msgChannel, ok := model.OnlineUsers[errorEvent.Username]
		mu.Unlock()

		if ok {
			msgChannel <- model.Message{
				MsgType: "logs",
				Msg: model.Logs{
					EventTime:      errorEvent.EventTime,
					ErrorCode:      errorEvent.ErrorCode,
					ErrorMessage:   errorEvent.ErrorMessage,
					ServiceName:    errorEvent.ServiceName,
					EventName:      errorEvent.EventName,
					Username:       errorEvent.Username,
					Explanation:    errorExplanation,
					ErrorExplained: true,
				}}
		}

		// save to database
		err = db.SaveErrorExplanations(ctx, &errorEvent, errorExplanation)
		if err != nil {
			fmt.Printf("Error saving error explanation")
			continue
		}
		fmt.Println("success")
		time.Sleep(1000000000) // for simulation
	}

}

