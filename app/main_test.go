package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type TestStatus struct {
	Name   string   `json:"name"`
	Passed bool     `json:"passed"`
	Errors []string `json:"errors"`
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/add", AddHandler)
	return router
}

func TestAddHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid addition",
			query:          "/add?num=1&num=20",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"success":true,"sum":21}`,
		},
		{
			name:           "Invalid number param",
			query:          "/add?num=1&num=x",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid number: x","success":false}`,
		},
		{
			name:           "Missing query param",
			query:          "/add",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"missing query param 'num'","success":false}`,
		},
	}

	r := SetUpRouter()

	testMap := make([]TestStatus, 0, len(tests))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			currStatus := &TestStatus{Name: tt.name}
			currStatus.Errors = []string{}

			if w.Code != tt.expectedStatus {
				errMsg := fmt.Sprintf("expected status %d, got %d", tt.expectedStatus, w.Code)
				currStatus.Errors = append(currStatus.Errors, errMsg)
				t.Error(errMsg)
			}
			if w.Body.String() != tt.expectedBody {
				errMsg := fmt.Sprintf("expected body %s, got %s", tt.expectedBody, w.Body.String())
				currStatus.Errors = append(currStatus.Errors, errMsg)
				t.Error(errMsg)
			}

			currStatus.Passed = len(currStatus.Errors) == 0

			testMap = append(testMap, *currStatus)
		})
	}

	authToken := os.Getenv("GCP_AUTH_TOKEN")
	if authToken == "" {
		panic("MISSING AUTH TOKEN")
	}

	cloudFnURL := os.Getenv("CLOUD_FN_URL")
	if cloudFnURL == "" {
		panic("MISSING CLOUD FUNCTION URL")
	}

	report := struct {
		Status  string       `json:"status"`
		Results []TestStatus `json:"results"`
	}{
		Status:  "completed",
		Results: testMap,
	}

	jsonData, jsonErr := json.Marshal(report)
	if jsonErr != nil {
		panic("FAILED TO MARSHAL JSON")
	}

	req, err := http.NewRequest("POST", cloudFnURL, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(fmt.Sprintf("FAILED TO CREATE REQUEST: %s", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to send results: %s", err))
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
	} else {
		fmt.Println("Function response:", string(bodyBytes))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic(fmt.Sprintf("Cloud function responded with error status: %d", resp.StatusCode))
	}
	fmt.Println("Sent test results to cloud function with status:", resp.StatusCode)
}
