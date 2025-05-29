package http_base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client  *http.Client
	BaseURL string
	Headers map[string]string
}

func New(base string, headers map[string]string) Client {

	client := &http.Client{Timeout: 15 * time.Second}

	return Client{
		client:  client,
		BaseURL: base,
		Headers: headers,
	}
}

func (c *Client) Get(url string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodGet, c.BaseURL+url, nil)

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Post(url string, payload []byte) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+url, bytes.NewBuffer(payload))

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Put(url string, payload []byte) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPut, c.BaseURL+url, bytes.NewBuffer(payload))

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Delete(url string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+url, nil)

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Path will safely join the base URL and the provided path and return a string
// that can be used in a request.
func (c *Client) Path(url string) string {

	base := strings.TrimRight(c.BaseURL, "/")

	if url == "" {
		return base
	}

	return base + "/" + strings.TrimLeft(url, "/")
}

// Pathf will call fmt.Sprintf with the provided values and then pass them
// to Client.Path as a convenience.
func (c *Client) Pathf(url string, v ...any) string {
	url = fmt.Sprintf(url, v...)
	return c.Path(url)
}

func Decode(r *http.Response, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}
