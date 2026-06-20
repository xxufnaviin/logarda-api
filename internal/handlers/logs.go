package handlers

import (
	"context"
	"fmt"
	"logarda/internal/db"
	"logarda/internal/model"
	"logarda/utils"
)

var errorMsg string
var err error
var errorEvent model.AWSErrorEvent

// @ goroutine
func ErrorLogsWorker() {
	ctx := context.Background()

	// consume error event from redis queue
	errorMsg, err = db.ConsumeErrorEvents()
	fmt.Println(errorMsg)
	if err != nil {
		fmt.Printf("Error getting event data.")
		return // change to continue
	}

	// make api call
	errorExplanation := "error explained" // placeholder

	// unmarshal string to JSON before saving to database
	err = utils.UnmarshalAWSErrorEvent(errorMsg, &errorEvent)
	if err != nil {
		fmt.Printf("Error during parsing event data.")
		return // change to continue
	}

	// stream to websocket (online users)

	// save to database
	err = db.SaveErrorExplanations(ctx, &errorEvent, errorExplanation)
	if err != nil {
		fmt.Printf("Error saving error explanation")
		return // change to continue
	}
	fmt.Println("success")

}
