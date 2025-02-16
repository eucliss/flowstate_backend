package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type OpenObserve struct {
	URL      string
	username string
	password string
	client   *http.Client
}

type OpenObserveInitializer struct{}

func (o *OpenObserveInitializer) Initialize() (Source, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("failed to load .env file: %v", err)
	}
	baseURL := os.Getenv("OPEN_OBSERVE_BASE_URL")
	username := os.Getenv("OPEN_OBSERVE_USERNAME")
	password := os.Getenv("OPEN_OBSERVE_PASSWORD")
	client := &http.Client{}

	return &OpenObserve{
		URL:      baseURL,
		username: username,
		password: password,
		client:   client,
	}, nil

}

func (o *OpenObserve) Query(query Query) []interface{} {
	endpoint := o.URL + "/api/" + query.SourceType + "/_search"

	initialQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"sql":        query.Query,
			"from":       0,
			"size":       query.Limit,
			"start_time": query.Start,
			"end_time":   query.End,
			"sql_mode":   "full",
		},
		"search_type": "ui",
	}
	jsonValue, err := json.Marshal(initialQuery)
	if err != nil {
		fmt.Printf("failed to marshal query: %v", err)
	}
	fmt.Printf("jsonValue: %s\n", string(jsonValue))

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("failed to create request: %v", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(o.username, o.password)

	resp, err := o.client.Do(req)
	if err != nil {
		fmt.Printf("failed to post logs: %v", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v", err)
		return nil
	}
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	fmt.Printf("result: %v\n", result)
	fmt.Println("--------------------------------")

	fmt.Printf("result hits: %v\n", result["hits"])
	fmt.Println("--------------------------------")
	return result["hits"].([]interface{})
}
