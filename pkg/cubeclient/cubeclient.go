package cubeclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SVCResponseOK struct {
	ReturnCode int32  `json:"return_code"`
	ClientID   string `json:"client_id"`
	ClientType int32  `json:"client_type"`
	Username   string `json:"username"`
	ExpiresIn  int32  `json:"expires_in"`
	UserID     int64  `json:"user_id"`
}

type SVCResponseERROR struct {
	ReturnCode  int32  `json:"return_code"`
	ErrorString string `json:"error_string"`
}

// CheckToken validates token with scope by
// sending request to the given token server.
// Server and client use CUBE OAUTH2 protocol
func CheckToken(
	host string,
	port int,
	token string,
	scope string) string {

	requestBodyFormat := `{
		"svc_msg": 1,
		"token": "%s",
		"scope": "%s"
	}`
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s:%v", host, port),
		bytes.NewBufferString(fmt.Sprintf(requestBodyFormat, token, scope))) //BTW check for error
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SVC-id", "2")
	req.Header.Set("Request-id", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		body := SVCResponseOK{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			panic(err)
		}
		prettyResponse, err := prettyResponseBody(body)
		if err != nil {
			panic(err)
		}
		return prettyResponse
	}
	body := SVCResponseERROR{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		panic(err)
	}
	prettyResponse, err := prettyResponseBody(body)
	if err != nil {
		panic(err)
	}
	return prettyResponse
}

func prettyResponseBody(respBody interface{}) (string, error) {
	switch respBody := respBody.(type) {
	default:
		return "", fmt.Errorf(
			"Unsupported response body type to pretty: %T",
			respBody)
	case SVCResponseOK:
		return fmt.Sprintf(
			"client_id: %s \n"+
				"client_type: %v \n"+
				"expires_in: %v \n"+
				"user_id: %v \n"+
				"username: %s \n",
			respBody.ClientID,
			respBody.ClientType,
			respBody.ExpiresIn,
			respBody.UserID,
			respBody.Username,
		), nil
	case SVCResponseERROR:
		return fmt.Sprintf(
			"error: %v \n"+
				"message: %s \n",
			respBody.ReturnCode,
			respBody.ErrorString,
		), nil
	}
}
