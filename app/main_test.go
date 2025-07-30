package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.query, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}
