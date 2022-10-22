package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	updatesMethod     = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host       string
	basePath   string
	clientHTTP http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:       host,
		basePath:   newBasePath(token),
		clientHTTP: http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(updatesMethod, q)
	if err != nil {
		return nil, fmt.Errorf("failed to do http-request %w", err)
	}

	var res UpdateResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json %w", err)
	}
	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("unable to send message %w", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) (bytes []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't do request %w", err)
	}
	req.URL.RawQuery = query.Encode()

	res, err := c.clientHTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http-request to foreign API failed %w", err)
	}
	defer func(f func() error) {
		errClose := f()
		if err == nil {
			err = errClose
		} else if errClose != nil {
			log.Printf("can't close request body %w", err)
		}
	}(res.Body.Close)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read http-response body %w", err)
	}

	return body, nil
}
