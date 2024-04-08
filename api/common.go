package api

import (
	"errors"
	"net/http"
)

var (
	ErrBadToken = errors.New("bad token (unauthorized)")
)

func VerifyToken(url string, token string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrBadToken
	}
	return nil
}
