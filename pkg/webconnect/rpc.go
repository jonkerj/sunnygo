package webconnect

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type ValueRequest struct {
	DestDev []string `json:"destDev"`
}

type AuthRequest struct {
	Right    string `json:"right"`
	Password string `json:"pass"`
}

type AuthResponse struct {
	Result struct {
		Sid      *string `json:"sid"`
		LoggedIn *bool   `json:"isLogin"`
	} `json:"result"`
}

type ErrorResponse struct {
	Code int `json:"err"`
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

func NewWebconnect(smaURL string) (*WebConnect, error) {
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

func (w *WebConnect) http(uri string, method string, body *bytes.Buffer, response interface{}) error {
	u := *w.url
	u.Path = uri

	if w.session != nil {
		q := u.Query()
		q.Set("sid", *w.session)
		u.RawQuery = q.Encode()
	}

	// NewRequest checks the _type_ of body, not if it's nil
	var req *http.Request
	if body == nil { // pass a nil-typed nil as body
		r, err := http.NewRequest(method, u.String(), nil)
		if err != nil {
			return fmt.Errorf("error creating http request: %v", err)
		}
		req = r
	} else {
		r, err := http.NewRequest(method, u.String(), body)
		if err != nil {
			return fmt.Errorf("error creating http request: %v", err)
		}
		req = r
	}

	if method != "GET" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error doing http request: %v", err)
	}
	defer resp.Body.Close()

	responsBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	errorResp := &ErrorResponse{}
	err = json.Unmarshal(responsBytes, &errorResp)
	if err == nil && errorResp.Code > 0 { // see if we can de-marshall into error
		return fmt.Errorf("API error: %d", errorResp.Code)
	}

	err = json.Unmarshal(responsBytes, response)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	return nil
}

func (w *WebConnect) get(uri string, response interface{}) error {
	return w.http(uri, "GET", nil, response)
}

func (w *WebConnect) post(uri string, data interface{}, response interface{}) error {
	body := bytes.NewBufferString("{}")
	if data != nil {
		bodyData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshalling to json: %v", err)
		}
		body = bytes.NewBuffer(bodyData)
	}

	return w.http(uri, "POST", body, response)
}

func (w *WebConnect) CheckSession() (bool, error) {
	if w.session == nil {
		return false, nil
	}

	resp := SessionCheckResponse{}
	err := w.post("/dyn/sessionCheck.json", nil, &resp)
	if err != nil {
		return false, fmt.Errorf("error checking session: %v", err)
	}

	if resp.Result.FreeSessions > 0 {
		return true, nil
	}

	return false, nil
}

func (w *WebConnect) DownloadLanguage() (*Language, error) {
	lStrings := map[string]string{}
	err := w.get("/data/l10n/en-US.json", &lStrings)
	if err != nil {
		return nil, fmt.Errorf("error downloading language: %w", err)
	}

	language := &Language{}

	for idStr, translation := range lStrings {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			continue // skip this entry
		}
		(*language)[idInt] = translation
	}

	return language, nil
}

func (w *WebConnect) DownloadMeta() (*Meta, error) {
	model, err := w.DownloadModel()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model: %w", err)
	}

	language, err := w.DownloadLanguage()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve language: %w", err)
	}
	meta := &Meta{
		model:    model,
		language: language,
	}

	return meta, nil
}

func (w *WebConnect) DownloadModel() (*Model, error) {
	model := &Model{}

	err := w.get("/data/ObjectMetadata_Istl.json", &model)
	if err != nil {
		return nil, fmt.Errorf("error downloading model: %w", err)
	}

	return model, nil
}

func (w *WebConnect) DownloadValues() (*ResultReponse, error) {
	values := &ResultReponse{}

	data := ValueRequest{
		DestDev: []string{},
	}
	err := w.post("/dyn/getAllOnlValues.json", data, &values)
	if err != nil {
		return nil, fmt.Errorf("error downloading values: %w", err)
	}

	return values, nil
}

func (w *WebConnect) Login(right, password string) error {
	data := AuthRequest{
		Right:    right,
		Password: password,
	}

	resp := AuthResponse{}
	err := w.post("/dyn/login.json", &data, &resp)
	if err != nil {
		return fmt.Errorf("error logging in: %v", err)
	}

	if resp.Result.Sid == nil {
		return errors.New("could not get session ID from API")
	}

	w.session = resp.Result.Sid

	return nil
}

func (w *WebConnect) Logout() error {
	if w.session != nil {
		resp := AuthResponse{}
		err := w.post("/dyn/logout.json", nil, &resp)
		if err != nil {
			return fmt.Errorf("error logging out: %v", err)
		}

		if resp.Result.LoggedIn == nil || *(resp.Result.LoggedIn) != false {
			return fmt.Errorf("logging out did not work")
		}
	}

	return nil
}
