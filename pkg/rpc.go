package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AuthResponse struct {
	Result struct {
		Sid      *string `json:"sid"`
		LoggedIn *bool   `json:"isLogin"`
	} `json:"result"`
}

type SessionCheckResponse struct {
	Result struct {
		FreeSessions int `json:"cntFreeSess"`
	} `json:"result"`
}

type ValueList []*json.RawMessage
type SubFields map[string]ValueList
type Fields map[string]SubFields
type ResultReponse struct {
	Result map[string]Fields `json:"result"`
}

type Nodifyable interface {
	Nodify(string, *Meta) (Node, error)
}

type IntValue struct {
	Val int `json:"val"`
}

type DurationValue struct {
	Val int `json:"val"`
}

type StringValue struct {
	Val string `json:"val"`
}

type TagListValue struct {
	Val []struct {
		Tag int `json:"tag"`
	}
}

type WebConnect struct {
	session    *string
	url        *url.URL
	httpClient *http.Client
}

func NewWebconnect(smaURL, right, password string) (*WebConnect, error) {
	c := WebConnect{session: nil, url: nil}

	u, err := url.Parse(smaURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c.httpClient = &http.Client{Transport: tr}

	c.url = u

	return &c, nil
}

func (w *WebConnect) http(uri string, method string, data *map[string]string, response interface{}) error {
	u := *w.url
	u.Path = uri

	if w.session != nil {
		q := u.Query()
		q.Set("sid", *w.session)
		u.RawQuery = q.Encode()
	}

	var body *bytes.Buffer = nil
	if method != "GET" {
		body = bytes.NewBufferString("{}")
		if data != nil {
			bodyData, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("error marshalling to json: %v", err)
			}
			body = bytes.NewBuffer(bodyData)
		}
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	if method != "GET" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	fmt.Printf("req: %v\n", req)

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error doing http request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("resp: %v\n", resp)
	responsBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("data: %v\n", string(responsBytes))
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	err = json.Unmarshal(responsBytes, response)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	return nil
}

func (s *WebConnect) Login(right, password string) error {
	var data = map[string]string{
		"right": right,
		"pass":  password,
	}

	resp := AuthResponse{}
	err := s.http("/dyn/login.json", "POST", &data, &resp)
	if err != nil {
		return fmt.Errorf("error logging in: %v", err)
	}

	fmt.Printf("%v\n", resp.Result)

	if resp.Result.Sid == nil {
		return errors.New("could not get session ID from API")
	}

	s.session = resp.Result.Sid

	return nil
}

func (s *WebConnect) Logout() error {
	if s.session != nil {
		fmt.Printf("Logging out")
		resp := AuthResponse{}
		err := s.http("/dyn/logout.json", "POST", nil, &resp)
		if err != nil {
			fmt.Printf("error logging out: %v\n", err)
			return fmt.Errorf("error logging out: %v", err)
		}

		if resp.Result.LoggedIn == nil || *(resp.Result.LoggedIn) != false {
			return fmt.Errorf("logging out did not work")
		}
		fmt.Printf("Logged out")
	}

	return nil
}

func (s *WebConnect) CheckSession() (bool, error) {
	if s.session == nil {
		return false, nil
	}

	resp := SessionCheckResponse{}
	err := s.http("/dyn/sessionCheck.json", "POST", nil, &resp)
	if err != nil {
		return false, fmt.Errorf("error checking session: %v", err)
	}

	if resp.Result.FreeSessions > 0 {
		return true, nil
	}

	return false, nil
}
