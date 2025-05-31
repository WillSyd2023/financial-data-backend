package util

import (
	"encoding/json"
	"io"
	"net/http"
)

type HttpClientItf interface {
	Get(string) (*http.Response, error)
	ReadAll(io.Reader) ([]byte, error)
	Unmarshal([]byte, any) error
}

type HttpClient struct {
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (hc *HttpClient) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}

func (hc *HttpClient) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func (hc *HttpClient) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
