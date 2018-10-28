package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	igc "github.com/marni/goigc"
)

func Test_postHandler_NotFound(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(Handler2))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotFound %d, received %d. ", 400, resp.StatusCode)
		return
	}

}

func Test_postHandler_Delete_Method(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(postHANDLER1))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotFound %d, received %d. ", 400, resp.StatusCode)
		return
	}

}

func Test_getHandler_Delete_Method(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getHANDLER1))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotFound %d, received %d. ", 400, resp.StatusCode)
		return
	}

}

func Test_Handler2_Delete_Method(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(Handler2))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotFound %d, received %d. ", 400, resp.StatusCode)
		return
	}

}

func Test_postHandler_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(postHANDLER1))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/randomTextJudtForTesting/",
		ts.URL + "/randommm/1sd5asd/",
		ts.URL + "/test",
		ts.URL + "/test/44a",
		ts.URL + "/45/45",
		ts.URL + "/12/asd",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}

func Test_getHandler_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getHANDLER1))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/randomTextJudtForTesting/",
		ts.URL + "/randommm/1sd5asd/",
		ts.URL + "/test",
		ts.URL + "/test/44a",
		ts.URL + "/45/45",
		ts.URL + "/12/asd",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}

func Test_Handler0_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Handler))
	defer ts.Close()

	testCases := []string{
		ts.URL + "/randomTextJudtForTesting",
		ts.URL + "/randommm/1sd5asd/",
		ts.URL + "/test",
		ts.URL + "/test44a",
		ts.URL + "/45as45",
		ts.URL + "/12da",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 400, resp.StatusCode)
			return
		}
	}
}

func Test_Handler2_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Handler2))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/randomTextJudtForTesting",
		ts.URL + "/randommm/1sd5asd/",
		ts.URL + "/test",
		ts.URL + "/test44a",
		ts.URL + "/45as45",
		ts.URL + "/12da",
		ts.URL + "/12475",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 400, resp.StatusCode)
			return
		}
	}
}

func Test_Handler3_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Handler3))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)

		if resp.StatusCode != 400 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", ts.URL, 400, resp.StatusCode)
			return
		}
	}
}

func Test_Handler6_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Handler6))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/randomTextJudtForTesting",
		ts.URL + "/randommm/1sd5asd",
		ts.URL + "/test",
		ts.URL + "/test44a",
		ts.URL + "/45as45",
		ts.URL + "/12da",
		ts.URL + "/1221312",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 400, resp.StatusCode)
			return
		}
	}
}

func Test_mongoConnect(t *testing.T) {
	if conn := connectDB("webhooks"); conn == nil {
		t.Error("No connection")
	}
}

func Test_validateURL(t *testing.T) {
	urlExists := validateURL(connectDB("igcTracks"), "random url", "url")
	if urlExists != 0 {
		t.Error("Track should not exist")
	}
}

func Test_returnTracks(t *testing.T) {

	result := returnTracks(1, 1)

	if result == "" {
		t.Error("Name should not exist")
	}
}

func Test_trackLength(t *testing.T) {

	result := trackLength(igc.Track{})

	if result != 0.0 {
		t.Error("Track was empty so it should have been 0.0!")
	}
}

func Test_insertToDB(t *testing.T) {

	insertToDB(collection, trackDB{})

}

func Test_respHandler6(t *testing.T) {

	resp, j := respHandler6("random")

	if resp == "" || j != 0 {
		t.Error("TimeStamp was random but it should at least return a json body")
	}

}

func Test_elapsedTime(t *testing.T) {

	result := elapsedTime(time.Now())
	if result == "" {
		t.Error("The string can never be empty!")
	}
}
func Test_getJ(t *testing.T) {
	//conn := connectDB("igcTracks")
	result := getJ(collection, "random")

	if result != 0 {
		t.Error("Thi timestamp should not exist!")
	}
}

func Test_whPerios(t *testing.T) {
	//conn := connectDB("igcTracks")
	triggerWebhookPeriod()
}
func Test_triggerWebhook(t *testing.T) {
	//conn := connectDB("igcTracks")
	triggerWebhook(nil)
}
func Test_getHANDLER1(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/igcinfo/api/igc", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getHANDLER1)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusOK { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}
func Test_Handler2(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/igcinfo/api/igc/asd", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler2)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusBadRequest { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusBadRequest)
	}
}
func Test_Handler3(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/igcinfo/api/igc/0/pilot", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler3)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusBadRequest { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusBadRequest)
	}
}
func Test_Handler4(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/igcinfo/api/ticker/latest", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler4)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusOK { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}
func Test_Handler5(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/igcinfo/api/ticker", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler5)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusOK { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}
func Test_Handler6(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/igcinfo/api/ticker/random", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler6)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusBadRequest { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusBadRequest)
	}
}

func Test_WebHookHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/paragliding/api/webhook/new_track", nil)
	if err != nil {
		t.Error(err)
	}

	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getWebHookHandler)

	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusBadRequest { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusBadRequest)
	}
}

func Test_getWebHookHandler_Delete_Method(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getWebHookHandler))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotFound %d, received %d. ", 400, resp.StatusCode)
		return
	}

}
