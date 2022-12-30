package handlers

import (
	"forum/s"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleError(t *testing.T) {
	// Re-declare and initialise
	//Temp := template.Must(template.ParseGlob("../../frontend/static/*.html"))
	HTTPData := s.StatusData{
		StatusCode: 404,
		StatusMsg:  "Test error message",
	}

	// Create a sample HTTP request
	request, errTest := http.NewRequest("GET", "/", nil)
	if errTest != nil {
		t.Fatal(errTest)
	}

	// Create a response recorder to record the response
	responseRec := httptest.NewRecorder()

	// Use the error handler that we want to test
	Error(responseRec, request, HTTPData)

	// Check if the correct error message was written to the response
	if responseRec.Code != HTTPData.StatusCode {
		t.Errorf("\nexpected status code: ' %v '. \nGot: '%v'",
			HTTPData.StatusCode, responseRec.Code)
	}
}
