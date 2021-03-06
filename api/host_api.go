// Copyright 2013 The go-cloudflare AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/url"
)

type HostAPI struct {
	client  *Client
	baseURL *url.URL
}
