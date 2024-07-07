package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func RestApiCall(baseURL string, method string, body map[string]any, query map[string]any) ([]byte, error) {
	// Create a URL object
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Add query parameters to the URL
	queryParams := reqURL.Query()
	for key, value := range query {
		queryParams.Set(key, fmt.Sprintf("%v", value))
	}
	reqURL.RawQuery = queryParams.Encode()

	// Create a new HTTP client
	client := &http.Client{}

	// Marshal the body parameters to JSON (if provided)
	var reqBody []byte
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %v", err)
		}
	}

	// Create a new HTTP request with the specified method
	req, err := http.NewRequest(method, reqURL.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set appropriate headers, e.g., Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return respBody, nil
}
