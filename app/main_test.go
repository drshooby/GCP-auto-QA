package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	run "cloud.google.com/go/run/apiv2"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
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

func getAuthToken(ctx context.Context, audience string) (string, error) {
	ts, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return "", fmt.Errorf("failed to create ID token source: %w", err)
	}
	token, err := ts.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get ID token: %w", err)
	}
	return token.AccessToken, nil
}

func getCloudFunctionURL(ctx context.Context, projectID, region, functionName string) (string, error) {
	client, err := run.NewServicesClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	req := &runpb.GetServiceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/services/%s", projectID, region, functionName),
	}

	service, err := client.GetService(ctx, req)
	if err != nil {
		return "", err
	}

	return service.Uri, nil
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

	ctx := context.Background()

	cloudFnURL, err := getCloudFunctionURL(ctx, "automated-qa-runner", "us-west1", "auto-qa-runner-function")
	if err != nil {
		panic(fmt.Sprintf("Failed to get function URL: %v", err))
	}

	authToken, err := getAuthToken(ctx, cloudFnURL)

	if err != nil {
		fmt.Printf("Auth token error: %v\n", err)
		panic(fmt.Sprintf("FAILED TO GET AUTH TOKEN: %v", err))
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
