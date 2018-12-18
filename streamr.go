// Package streamr implements Go client for Streamr API.
package streamr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is the Streamr API client.
type Client struct {
	// HTTP Client used to communicate with the API.
	client *http.Client

	// APIKey is Streamr API key.
	APIKey string

	// Base URL for API requests.
	BaseURL *url.URL

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Streamr API.
	Data *DataService
}

type service struct {
	client *Client
}

// Response is a Streamr API response. This wraps standard http.Response.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// NewClient creates a new Streamr client.
func NewClient(apiKey string) (*Client, error) {
	var url = "https://www.streamr.com/api/v1/"
	return NewClientWithBaseURL(apiKey, url)
}

// NewClientWithBaseURL creates a new Streamr client with given baseURL for mostly testing purposes.
func NewClientWithBaseURL(apiKey, baseURL string) (*Client, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: 0,
		},
		APIKey:  apiKey,
		BaseURL: url,
	}
	c.common.client = c
	c.Data = (*DataService)(&c.common)
	return c, nil
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err = enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Close = true
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %v", c.APIKey))

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := res.Body.Close(); e != nil {
			err = e
		}
	}()
	response := newResponse(res)

	err = CheckResponse(res)
	if err != nil {
		return response, err
	}
	if v != nil {
		decErr := json.NewDecoder(res.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return response, err
}

// An ErrorResponse reports error caused by an API request.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
	)
}

// CheckResponse checks the API response for errors.
func CheckResponse(r *http.Response) error {
	if 200 <= r.StatusCode && r.StatusCode <= 299 {
		return nil
	}
	return &ErrorResponse{
		Response: r,
	}
}
