package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type SnykOptions struct {
	GroupId   string
	ApiKey    string
	UserAgent string
}

const rest_version = "2023-09-11"

var ErrInvalidAuthn = errors.New("credentials not valid")
var ErrInvalidAuthz = errors.New("credentials not authorized to access resource")
var ErrNotFound = errors.New("requested resource not found")
var ErrUnexpectedStatus = errors.New("unexpected HTTP status code")

func clientDo(so SnykOptions, method string, path string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, constructUrl(path), bytes.NewReader(body))

	generateHeaders(so, req)

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode < 300 {
		return res, nil
	} else if res.StatusCode == 401 {
		return nil, fmt.Errorf("%w", ErrInvalidAuthn)
	} else if res.StatusCode == 403 {
		return nil, fmt.Errorf("%w", ErrInvalidAuthz)
	} else if res.StatusCode == 404 {
		return nil, fmt.Errorf("%w", ErrNotFound)
	} else {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}
}

func clientDoRest(so SnykOptions, method string, path string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, constructUrlRest(path), bytes.NewReader(body))

	query := req.URL.Query()
	query.Set("version", rest_version)
	req.URL.RawQuery = query.Encode()

	generateHeaders(so, req)

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode < 300 {
		return res, nil
	} else if res.StatusCode == 401 {
		return nil, fmt.Errorf("%w", ErrInvalidAuthn)
	} else if res.StatusCode == 403 {
		return nil, fmt.Errorf("%w", ErrInvalidAuthz)
	} else if res.StatusCode == 404 {
		return nil, fmt.Errorf("%w", ErrNotFound)
	} else {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}
}
func generateHeaders(so SnykOptions, req *http.Request) {
	authToken := fmt.Sprintf("token %s", so.ApiKey)
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", so.UserAgent)
}
func constructUrl(path string) string {
	snykEndpoint := "https://snyk.io/api/v1%s"
	return fmt.Sprintf(snykEndpoint, path)
}

func constructUrlRest(path string) string {
	snykEndpoint := "https://api.snyk.io/rest%s"
	return fmt.Sprintf(snykEndpoint, path)
}
