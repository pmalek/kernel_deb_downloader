package http

import (
	"bytes"
	"io"
	"net/http"
)

// Getter interface provides one method: Get, to
// allow using it in places where we would like test
// an HTTP client
type Getter interface {
	Get(string) (*http.Response, error)
}

// Header interface provides one method: Head, to
// allow using it in places where we would like test
// an HTTP client
type Header interface {
	Head(string) (*http.Response, error)
}

// GetterHeader interface combines Header and Getter
// interfaces
type GetterHeader interface {
	Getter
	Header
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// MockedClient fulfills the Getter interface
type MockedClient struct {
	response   string
	err        error
	statusCode int
}

// SetResponse sets the response that MockedClient
// will receive when calling GET
func (c *MockedClient) SetResponse(response string) {
	c.response = response
}

// SetError sets the error that MockedClient will
// receive after calling GET
func (c *MockedClient) SetError(err error) {
	c.err = err
}

// SetStatusCode sets the statusCode on MockedClient
func (c *MockedClient) SetStatusCode(statusCode int) {
	c.statusCode = statusCode
}

// Get returns a preset response body and a preset error
func (c MockedClient) Get(url string) (*http.Response, error) {
	resp := &http.Response{
		Body:       nopCloser{bytes.NewBufferString(c.response)},
		StatusCode: c.statusCode,
	}

	return resp, c.err
}

// Head returns a preset response body and a preset error
func (c MockedClient) Head(url string) (*http.Response, error) {
	resp := &http.Response{
		Body:          nopCloser{bytes.NewBufferString(c.response)},
		StatusCode:    c.statusCode,
		ContentLength: int64(len(c.response)),
	}

	return resp, c.err
}
