package api

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	resty *resty.Client
}

func New(token string) *Client {
	client := resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("X-Access-Token", token)

	return &Client{
		resty: client,
	}
}
