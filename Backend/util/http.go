package util

import (
	"io"
	"net/http"
)

type HttpClientItf interface {
	Get(string) (*http.Response, error)
	ReadAll(io.Reader) ([]byte, error)
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
