package api

import (
	"errors"
	"io"
	"net/http"
)

var (
	ErrBadToken = errors.New("bad token (unauthorized)")
)

func VerifyToken(url string, token string) error {
	req, err := newRequest(http.MethodGet, url, token, nil)
	if err != nil {
		return err
	}

	resp, err := doRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrBadToken
	}
	return nil
}

func newRequest(method string, url string, token string, body io.Reader) (*http.Request, error) {
	req, err := newRequestNoToken(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	return req, nil
}

func newRequestNoToken(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func doRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
