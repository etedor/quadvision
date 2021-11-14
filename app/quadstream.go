package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const qsAPI = "https://quadstream.tv/stream/api/"

type qsCredentials struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type qsResponse struct {
	ShortID string `json:"short_id"`
}

type qsStreams struct {
	Stream1 string `json:"stream1"`
	Stream2 string `json:"stream2"`
	Stream3 string `json:"stream3"`
	Stream4 string `json:"stream4"`
}

func login(username, secret string) ([]*http.Cookie, *string, error) {
	u := qsAPI + "login"

	payload := &qsCredentials{
		Username: username,
		Secret:   secret,
	}
	p := new(bytes.Buffer)
	json.NewEncoder(p).Encode(payload)

	resp, err := post([]*http.Cookie{}, p, u)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var val qsResponse
	if err := json.Unmarshal([]byte(body), &val); err != nil {
		return nil, nil, err
	}

	return resp.Cookies(), &val.ShortID, nil
}

// post issues a POST to u
func post(cc []*http.Cookie, payload *bytes.Buffer, u string) (*http.Response, error) {
	req, err := http.NewRequest("POST", u, payload)
	if err != nil {
		return nil, err
	}

	for i := range cc {
		req.AddCookie(cc[i])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	return resp, nil
}

func update(cc []*http.Cookie, id string, ss []string) error {
	u := qsAPI + fmt.Sprintf("stream/%s/update", id)

	payload := &qsStreams{
		Stream1: ss[0],
		Stream2: ss[1],
		Stream3: ss[2],
		Stream4: ss[3],
	}
	p := new(bytes.Buffer)
	json.NewEncoder(p).Encode(payload)

	resp, err := post(cc, p, u)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
