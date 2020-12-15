package test

import (
	"fmt"
	"io"
	"net/http"
)

type TestRequest struct {
	Client      *http.Client
	Method      string
	URL         string
	ContentType string
	Token       string
	Body        io.Reader
}

func (tr TestRequest) Make() (*http.Response, error) {
	req, err := http.NewRequest(tr.Method, tr.URL, tr.Body)
	if err != nil {
		return nil, err
	}
	if tr.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tr.Token))
	}
	if tr.ContentType != "" {
		req.Header.Set("Content-Type", tr.ContentType)
	}
	return tr.Client.Do(req)
}
