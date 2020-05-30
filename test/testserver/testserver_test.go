package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testBody = `{
	"svc_msg": 1
	"token": "test1"
	"scope": "read"
}
`

func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	s := string(b)
	return s
}
func TestServer(t *testing.T) {
	srv := httptest.NewServer(Handlers())
	defer srv.Close()

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/checktoken", srv.URL), bytes.NewBufferString(testBody)) //BTW check for error
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SVC-id", "2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Printf("response Body: %s", string(body))
	/* res, err := http.Get(fmt.Sprintf("%s/rate/btc", srv.URL))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "BitCoin to USD rate: 0.000000 $\n" {
		t.Fail()
	} */
}
