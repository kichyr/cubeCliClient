package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const testBody = `{
	"svc_msg": 1,
	"token": "test1",
	"scope": "read"
}
`

func TestServer(t *testing.T) {
	srv := httptest.NewServer(Handlers())
	defer srv.Close()

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/checktoken", srv.URL), bytes.NewBufferString(testBody)) //BTW check for error
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SVC-id", "2")
	req.Header.Set("Request-id", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body := SVCResponseOK{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Printf("response Body: %v", (body))

	expectedBody := SVCResponseOK{1, "test_client_id", 2002, "testuser0@mail.ru", 3600, 0}
	if !reflect.DeepEqual(expectedBody, body) {
		t.Fatal(fmt.Sprintf("unexpected response body: %v, get: %v", body, expectedBody))
	}
}
