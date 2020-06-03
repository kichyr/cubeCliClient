package cubeclient

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kichyr/cubeCliClient/test/testserver"
)

const expectedAnswerOK = `client_id: test_client_id 
client_type: 2002 
expires_in: 3600 
user_id: 0 
username: testuser0@mail.ru 
`

const expectedAnswerBadScope = `error: 6 
message: CUBE_OAUTH2_ERR_BAD_SCOPE 
`

func TestClientOKResponse(t *testing.T) {
	serv := testserver.NewTestServer("../../test/testserver/test_tokens.json")
	httpTestSrv := httptest.NewServer(serv.Handlers())
	defer httpTestSrv.Close()

	url := strings.Split(httpTestSrv.URL, ":") // ["http", "//127.0.0.1", "port"]
	port, _ := strconv.Atoi(url[2])
	answer, _ := CheckToken(url[1][2:], port, "test1", "read")
	if answer != expectedAnswerOK {
		t.Fatal(fmt.Sprintf(
			"unexpected answer: %s, \n expected: %s",
			answer,
			expectedAnswerOK))
	}
}

func TestClientBadScope(t *testing.T) {
	serv := testserver.NewTestServer("../../test/testserver/test_tokens.json")
	httpTestSrv := httptest.NewServer(serv.Handlers())
	defer httpTestSrv.Close()

	url := strings.Split(httpTestSrv.URL, ":") // ["http", "//127.0.0.1", "port"]
	port, _ := strconv.Atoi(url[2])
	answer, _ := CheckToken(url[1][2:], port, "test1", "admin")
	fmt.Print("++++" + answer + "+++")
	if answer != expectedAnswerBadScope {
		t.Fatal(fmt.Sprintf(
			"unexpected answer: %s, \n expected: %s",
			answer,
			expectedAnswerBadScope))
	}
}
