package cubeclient

import (
	"fmt"
	"testing"
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
	answer1 := CheckToken("localhost", 8091, "test1", "read")
	if answer1 != expectedAnswerOK {
		t.Fatal(fmt.Sprintf(
			"unexpected answer: %s, \n expected: %s",
			answer1,
			expectedAnswerOK))
	}
}

func TestClientBadScope(t *testing.T) {
	answer := CheckToken("localhost", 8091, "test1", "admin")
	fmt.Print("++++" + answer + "+++")
	if answer != expectedAnswerBadScope {
		t.Fatal(fmt.Sprintf(
			"unexpected answer: %s, \n expected: %s",
			answer,
			expectedAnswerBadScope))
	}
}
