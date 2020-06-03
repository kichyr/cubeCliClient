package testserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const testOKBody = `{
	"svc_msg": 1,
	"token": "test1",
	"scope": "read"
}
`

const testBadScopeBody = `{
	"svc_msg": 1,
	"token": "test1",
	"scope": "admin"
}
`

func TestServerOKRespose(t *testing.T) {
	serv := NewTestServer("./test_tokens.json")
	srv := httptest.NewServer(serv.Handlers())
	defer srv.Close()

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/", srv.URL), bytes.NewBufferString(testOKBody)) //BTW check for error
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
	_ = json.NewDecoder(resp.Body).Decode(&body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Printf("response Body: %v", (body))

	expectedBody := SVCResponseOK{0, "test_client_id", 2002, "testuser0@mail.ru", 3600, 0}
	if !reflect.DeepEqual(expectedBody, body) {
		t.Fatal(fmt.Sprintf("unexpected response body: %v, expected: %v", body, expectedBody))
	}
}

func TestServerWrongToken(t *testing.T) {
	serv := NewTestServer("./test_tokens.json")
	srv := httptest.NewServer(serv.Handlers())
	defer srv.Close()

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(testBadScopeBody)) //BTW check for error
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SVC-id", "2")
	req.Header.Set("Request-id", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body := SVCResponseERROR{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Printf("response Body: %v", (body))

	expectedBody := SVCResponseERROR{6, "CUBE_OAUTH2_ERR_BAD_SCOPE"}
	if !reflect.DeepEqual(expectedBody, body) {
		t.Fatal(fmt.Sprintf("unexpected response body: %v, get: %v", body, expectedBody))
	}
}
