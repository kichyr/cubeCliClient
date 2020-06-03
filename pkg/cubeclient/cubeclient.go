package cubeclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
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

type SVCRequestBody struct {
	SVCMsg int32  `json:"svc_msg"`
	Token  string `json:"token"`
	Scope  string `json:"scope"`
}

// CheckToken validates token & scope by
// sending request to the given token server.
// Server and client use CUBE OAUTH2 protocol.
func CheckToken(
	host string,
	port int,
	token string,
	scope string) (string, error) {

	reqBody := SVCRequestBody{
		SVCMsg: 1,
		Token:  token,
		Scope:  scope,
	}

	reqBytesBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s:%v", host, port),
		bytes.NewBuffer(reqBytesBody),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SVC-id", "2")
	req.Header.Set("Request-id", fmt.Sprint(rand.Intn(1000)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error during executing request, error: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		body := SVCResponseOK{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return "", fmt.Errorf("Bad ok response body.")
		}
		prettyResponse, err := prettyResponseBody(body)
		if err != nil {
			return "", fmt.Errorf(
				"Unable to prettify body: %v, error: %s",
				body, err.Error())
		}
		return prettyResponse, nil
	}
	body := SVCResponseERROR{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", fmt.Errorf("Bad error response body.")
	}
	prettyResponse, err := prettyResponseBody(body)
	if err != nil {
		return "", fmt.Errorf(
			"Unable to prettify error body: %v, error: %s",
			body, err.Error())
	}
	return prettyResponse, nil
}

// Genrates message in human readable format according to
// given response body. It returns error if respBody has
// unexpected format.
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
