package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StatusCode int32

const (
	CUBE_OAUTH2_ERR_OK StatusCode = iota + 1
	CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND
	CUBE_OAUTH2_ERR_DB_ERROR
	CUBE_OAUTH2_ERR_UNKNOWN_MSG
	CUBE_OAUTH2_ERR_BAD_PACKET
	CUBE_OAUTH2_ERR_BAD_CLIENT
	CUBE_OAUTH2_ERR_BAD_SCOPE
)

func (code StatusCode) String() string {
	return fmt.Sprintf("%v", int32(code))
}

type SVCResponseOK struct {
	ReturnCode int32  `json:"return_code"`
	ClientID   string `json:"client_id"`
	ClientType string `json:"client_type"`
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

type TokenStorage struct {
	tokensInfo []TokenInfo
}
type TokenInfo struct {
	UserID    int32    `json:"user_id"`
	Username  string   `json:"username"`
	ExpiresIn string   `json:"expires_in"`
	Token     string   `json:"token"`
	Scopes    []string `json:"scopes"`
}

func NewTokenStorage(path string) *TokenStorage {
	plan, _ := ioutil.ReadFile(path)
	ts := TokenStorage{make([]TokenInfo, 0)}
	err := json.Unmarshal(plan, &ts)
	if err != nil {
		panic("Cannot unmarshal the json.")
	}
	return &ts
}

type TokenChecker struct {
	tStorage TokenStorage
}

func NewTokenChecker(storagePath string) *TokenChecker {
	return &TokenChecker{*NewTokenStorage(storagePath)}
}

// checkToken check given token and scope according data in TokenStorage.
func (tc *TokenChecker) checkToken(token string, scope string) (*TokenInfo, error) {
	badScope := false
	for _, tokenInfo := range tc.tStorage.tokensInfo {
		if token == tokenInfo.Token {
			badScope = true
			for _, s := range tokenInfo.Scopes {
				if s == scope {
					return &tokenInfo, nil
				}
			}
		}
	}
	if badScope {
		return nil, fmt.Errorf(CUBE_OAUTH2_ERR_BAD_SCOPE.String())
	}

	return nil, fmt.Errorf(CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND.String())
}

func (ts TokenChecker) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	errorResponse := func(errorCode StatusCode, message string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		respBody, _ := json.Marshal(SVCResponseERROR{int32(errorCode), message})
		fmt.Fprint(w, respBody)
	}

	if req.Method != http.MethodPost {
		errorResponse(CUBE_OAUTH2_ERR_BAD_PACKET, "bad method")
		return
	}

	if contentType != "application/json" {
		errorResponse(CUBE_OAUTH2_ERR_BAD_PACKET, "unsupported content-type")
		return
	}

	svcID := req.Header.Get("SVC-id")
	requestID := req.Header.Get("Request-id")
	if svcID != "" || requestID == "" {
		errorResponse(CUBE_OAUTH2_ERR_BAD_PACKET, "required headers not specified")
		return
	}

	parsedBody := SVCRequestBody{}
	err := json.NewDecoder(req.Body).Decode(&parsedBody)

	if err != nil {
		errorResponse(CUBE_OAUTH2_ERR_BAD_PACKET, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = ts.checkToken(parsedBody.Token, parsedBody.Scope)
	if err != nil {
		if err.Error() == CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND.String() {
			errorResponse(CUBE_OAUTH2_ERR_BAD_PACKET, err.Error())
			return
		}

	}

}

func Handlers() http.Handler {
	r := http.NewServeMux()
	r.Handle("/checktoken", *NewTokenChecker("./test_tokens.json"))
	return r
}

func main() {

	http.Handle("/checktoken", *NewTokenChecker("./test_tokens.json"))
	http.ListenAndServe(":8091", nil)
}
