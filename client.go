package wdb

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

type Client struct {
	ServerURL string
}

func NewClient(serverURL string) *Client {
	return &Client{
		ServerURL: serverURL,
	}
}

func build_url(baseurl string, key string) string {
	return baseurl + "/" + key
}

func (c *Client) Set(key string, raw []byte) error {
	resp, err := http.Post(build_url(c.ServerURL, key), "", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New("set failure")
	}
	return nil
}

func (c *Client) Get(key string) ([]byte, error) {
	resp, err := http.Get(build_url(c.ServerURL, key))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("get failure")
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return raw, nil

}
