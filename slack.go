package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func sendSlackMessage(webhookURL, msg string) (bool, error) {
	payload := []byte(fmt.Sprintf(`{"text":"%s"}`, msg))
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}
