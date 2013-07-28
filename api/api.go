// Copyright 2013 The go-cloudflare AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
TODO: Add basic usage here
*/

package api

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1"
	userAgent      = "go-cloudflare/" + libraryVersion
	baseURLClient  = "https://www.cloudflare.com/api_json.html"
	baseURLHost    = "https://api.cloudflare.com/host-gw.html"
)

type Client struct {
	client    *http.Client
	config    *Config
	userAgent string

	ClientAPI *ClientAPI
	HostAPI   *HostAPI
}

type Config struct {
	Token   string
	Email   string
	HostKey string
	UserKey string
}

func NewClient(conf *Config) *Client {

	c := &Client{
		client:    http.DefaultClient,
		config:    conf,
		userAgent: userAgent,
	}

	baseURLForClient, _ := url.Parse(baseURLClient)
	baseURLForHost, _ := url.Parse(baseURLHost)

	c.ClientAPI = &ClientAPI{
		client:  c,
		baseURL: baseURLForClient,
	}
	c.HostAPI = &HostAPI{
		client:  c,
		baseURL: baseURLForHost,
	}
	return c
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}
