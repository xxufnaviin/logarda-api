package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"logarda/internal/config"
	"logarda/internal/model"
	"net/http"
)

var llmEndpoint = "analytics/llm/inference"

func GetErrorExplanation(w http.ResponseWriter, r *http.Request) {

}

func GetLLMInference(errorEvent *model.AWSErrorEvent) string {
	var response model.LLMInferenceResponse

	// convert the error event to bytes
	payload, err := json.Marshal(errorEvent)
	if err != nil {
		fmt.Println("Error converting payload. ", err)
		return ""
	}
	// create request
	requestUrl := fmt.Sprintf("%s%s", config.ANALYTICS_API, llmEndpoint)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request. ", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting response. ", err)
		return ""
	}
	defer resp.Body.Close()
	// fetch it and store it in response structure
	json.NewDecoder(resp.Body).Decode(&response)

	// concat the explanation and solution together as a string
	explanationBytes, _ := json.Marshal(response.Data)
	fmt.Println(response.Data)
	errorExplanation := string(explanationBytes)

	fmt.Print(errorExplanation)
	return errorExplanation
}
