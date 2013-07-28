// Copyright 2013 The go-cloudflare AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ClientAPI struct {
	client  *Client
	baseURL *url.URL
}

func (c *ClientAPI) NewRequest(data *url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", c.baseURL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

type StatsResponse struct {
	Response  *StatsResult `json:"response,omitempty"`
	Result    string       `json:"result,omitempty"`
	Message   string       `json:"msg,omitempty"`
	ErrorCode string       `json:"err_code,omitempty"`
}

type StatsResult struct {
	Result *Stats `json:"result,omitempty"`
}

type Stats struct {
	TimeZero int          `json:"timeZero,omitempty"`
	TimeEnd  int          `json:"timeEnd,omitempty"`
	Count    int          `json:"count,omitempty"`
	HasMore  bool         `json:"has_more,omitempty"`
	Objects  []StatObject `json:"objs,omitempty"`
}

type StatObject struct {
	CachedServerTime    int               `json:"cachedServerTime,omitempty"`
	CachedExpryTime     int               `json:"cachedExpryTime,omitempty"`
	TrafficBreakdown    *TrafficBreakdown `json:"trafficBreakdown,omitempty"`
	BandwidthServed     *BandwidthServed  `json:"bandwidthServed,omitempty"`
	RequestsServed      *RequestsServed   `json:"requestsServed,omitempty"`
	ProZone             bool              `json:"pro_zone,omitempty"`
	PageLoadTime        int               `json:"pageLoadTime,omitempty"`
	CurrentServerTime   int               `json:"currentServerTime,omitempty"`
	Interval            int               `json:"interval,omitempty"`
	ZoneCDate           int               `json:"zoneCDate,omitempty"`
	UserSecuritySetting string            `json:"userSecuritySetting,omitempty"`
	DevMode             int               `json:"dev_mode,omitempty"`
	IPv46               int               `json:"ipv46,omitempty"`
	OB                  int               `json:"ob,omitempty"`
	CacheLevel          string            `json:"cache_lvl,omitempty"`
}

type TrafficBreakdown struct {
	PageViews *TrafficStat `json:"pageviews,omitempty"`
	Uniques   *TrafficStat `json:"pageviews,omitempty"`
}

type TrafficStat struct {
	Regular int `json:"regular,omitempty"`
	Threat  int `json:"threat,omitempty"`
	Crawler int `json:"crawler,omitempty"`
}

type BandwidthServed struct {
	CloudFlare float64 `json:"cloudflare,omitempty"`
	User       float64 `json:"user,omitempty"`
}

type RequestsServed struct {
	CloudFlare int `json:"cloudflare,omitempty"`
	User       int `json:"user,omitempty"`
}

func (c *ClientAPI) Stats(z string, interval int) (*Stats, error) {
	data := &url.Values{
		"tkn":      {c.client.config.Token},
		"email":    {c.client.config.Email},
		"a":        {"stats"},
		"z":        {z},
		"interval": {strconv.Itoa(interval)},
	}
	req, err := c.NewRequest(data)
	if err != nil {
		return nil, err
	}

	v := new(StatsResponse)
	_, err = c.client.Do(req, v)
	if err != nil {
		return nil, err
	}
	if v.Result == "error" || v.ErrorCode != "" {
		return nil, fmt.Errorf("%v %v: %v", req.Method, req.URL, v.Message)
	}

	return v.Response.Result, nil
}

type ZoneLoadMultiResponse struct {
	Response  *ZoneLoadMultiZones `json:"response,omitempty"`
	Result    string              `json:"result,omitempty"`
	Message   string              `json:"msg,omitempty"`
	ErrorCode string              `json:"err_code,omitempty"`
}

type ZoneLoadMultiZones struct {
	Zones *Zones `json:"zones,omitempty"`
}

type Zones struct {
	HasMore bool   `json:"has_more,omitempty"`
	Count   int    `json:"count,omitempty"`
	Objects []Zone `json:"objs,omitempty"`
}

type Zone struct {
	ZoneID          string `json:"zone_id,omitempty"`
	UserID          string `json:"user_id,omitempty"`
	ZoneName        string `json:"zone_name,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	ZoneStatus      string `json:"zone_status,omitempty"`
	ZoneMode        string `json:"zone_mode,omitempty"`
	ZoneType        string `json:"zone_type,omitempty"`
	HostID          string `json:"host_id,omitempty"`
	HostPubName     string `json:"host_pubname,omitempty"`
	HostWebsite     string `json:"host_website,omitempty"`
	VTXT            string `json:"vtxt,omitempty"`
	FQDNS           string `json:"fqdns,omitempty"`
	Step            int    `json:"step,omitempty"`
	ZoneStatusClass string `json:"zone_status_class,omitempty"`
	ZoneStatusDesc  string `json:"zone_status_desc,omitempty"`
	// TODO uclear type of ns_vanity_map, orig_registrar, orig_dnshost, orig_ns_names
	Props       map[string]int    `json:"props,omitempty"`
	ConfirmCode map[string]string `json:"confirm_code,omitempty"`
	Allow       []string          `json:"allow,omitempty"`
}

func (c *ClientAPI) ZoneLoadMulti() (*Zones, error) {
	data := &url.Values{
		"tkn":   {c.client.config.Token},
		"email": {c.client.config.Email},
		"a":     {"zone_load_multi"},
	}
	req, err := c.NewRequest(data)
	if err != nil {
		return nil, err
	}

	v := new(ZoneLoadMultiResponse)
	_, err = c.client.Do(req, v)
	if err != nil {
		return nil, err
	}
	if v.Result == "error" || v.ErrorCode != "" {
		return nil, fmt.Errorf("%v %v: %v", req.Method, req.URL, v.Message)
	}

	return v.Response.Zones, nil
}

type ZoneCheckResponse struct {
	Response  *ZoneCheck `json:"response,omitempty"`
	Result    string     `json:"result,omitempty"`
	Message   string     `json:"msg,omitempty"`
	ErrorCode string     `json:"err_code,omitempty"`
}

type ZoneCheck struct {
	Zones map[string]int `json:"zones,omitempty"`
}

func (c *ClientAPI) ZoneCheck(zones []string) (map[string]int, error) {
	data := &url.Values{
		"tkn":   {c.client.config.Token},
		"email": {c.client.config.Email},
		"a":     {"zone_check"},
		"zones": {strings.Join(zones, ",")},
	}
	req, err := c.NewRequest(data)
	if err != nil {
		return nil, err
	}

	v := new(ZoneCheckResponse)
	_, err = c.client.Do(req, v)
	if err != nil {
		return nil, err
	}
	if v.Result == "error" || v.ErrorCode != "" {
		return nil, fmt.Errorf("%v %v: %v", req.Method, req.URL, v.Message)
	}

	return v.Response.Zones, nil
}

type ZoneIPsResponse struct {
	Response  *ZoneIPs `json:"response,omitempty"`
	Result    string   `json:"result,omitempty"`
	Message   string   `json:"msg,omitempty"`
	ErrorCode string   `json:"err_code,omitempty"`
}

type ZoneIPs struct {
	IPs []ZoneIP `json:"ips,omitempty"`
}

type ZoneIP struct {
	IP             string  `json:"ip,omitempty"`
	Classification string  `json:"classification,omitempty"`
	Hits           int     `json:"hits,omitempty"`
	Latitude       float64 `json:"latitude,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
	ZoneName       string  `json:"zone_name,omitempty"`
}

func (c *ClientAPI) ZoneIPs(z string, hours int, class string) ([]ZoneIP, error) {
	data := &url.Values{
		"tkn":   {c.client.config.Token},
		"email": {c.client.config.Email},
		"a":     {"zone_ips"},
		"z":     {z},
		"hours": {strconv.Itoa(hours)},
	}
	if class != "" {
		data.Add("class", class)
	}

	req, err := c.NewRequest(data)
	if err != nil {
		return nil, err
	}

	v := new(ZoneIPsResponse)
	_, err = c.client.Do(req, v)
	if err != nil {
		return nil, err
	}
	if v.Result == "error" || v.ErrorCode != "" {
		return nil, fmt.Errorf("%v %v: %v", req.Method, req.URL, v.Message)
	}

	return v.Response.IPs, nil
}

type ZoneSettingResponse struct {
	Response  *ZoneSettingResult `json:"response,omitempty"`
	Result    string             `json:"result,omitempty"`
	Message   string             `json:"msg,omitempty"`
	ErrorCode string             `json:"err_code,omitempty"`
}

type ZoneSettingResult struct {
	Result *ZoneSettingObjects `json:"result,omitempty"`
}

type ZoneSettingObjects struct {
	Objects []ZoneSetting `json:"objs,omitempty"`
}

type ZoneSetting struct {
	UserSecuritySetting string `json:"userSecuritySetting,omitempty"`
	DevMode             int    `json:"dev_mode,omitempty"`
	IPv46               int    `json:"ipv46,omitempty"`
	OB                  int    `json:"ob,omitempty"`
	CacheLevel          string `json:"cache_lvl,omitempty"`
	OutboundLinks       string `json:"outboundLinks,omitempty"`
	Async               string `json:"async,omitempty"`
	Bic                 string `json:"bic,omitempty"`
	CacheTTL            string `json:"chl_ttl,omitempty"`
	ExpireTTL           string `json:"exp_ttl,omitempty"`
	FPurgeTS            string `json:"fpurge_ts,omitempty"`
	Hotlink             string `json:"hotlink,omitempty"`
	IMG                 string `json:"img,omitempty"`
	Lazy                string `json:"lazy,omitempty"`
	Minify              string `json:"minify,omitempty"`
	Outlink             string `json:"outlink,omitempty"`
	Preload             string `json:"preload,omitempty"`
	S404                string `json:"s404,omitempty"`
	SecurityLevel       string `json:"sec_lvl,omitempty"`
	SPDY                string `json:"spdy,omitempty"`
	SSL                 string `json:"ssl,omitempty"`
	WAFProfile          string `json:"waf_profile,omitempty"`
}

func (c *ClientAPI) ZoneSettings(z string) ([]ZoneSetting, error) {
	data := &url.Values{
		"tkn":   {c.client.config.Token},
		"email": {c.client.config.Email},
		"a":     {"zone_settings"},
		"z":     {z},
	}

	req, err := c.NewRequest(data)
	if err != nil {
		return nil, err
	}

	v := new(ZoneSettingResponse)
	_, err = c.client.Do(req, v)
	if err != nil {
		return nil, err
	}
	if v.Result == "error" || v.ErrorCode != "" {
		return nil, fmt.Errorf("%v %v: %v", req.Method, req.URL, v.Message)
	}

	return v.Response.Result.Objects, nil
}

type RecLoadAllResponse struct {
	Response  *RecordsResponse `json:"response,omitempty"`
	Result    string           `json:"result,omitempty"`
	Message   string           `json:"msg,omitempty"`
	ErrorCode string           `json:"err_code,omitempty"`
}

type RecordsResponse struct {
	Records *Records `json:"recs,omitempty"`
}

type Records struct {
	HasMore bool     `json:"has_more,omitempty"`
	Count   int      `json:"count,omitempty"`
	Objects []Record `json:"objs,omitempty"`
}

type Record struct {
	RecordID       string          `json:"rec_id,omitempty"`
	RecordTag      string          `json:"rec_tag,omitempty"`
	ZoneName       string          `json:"zone_name,omitempty"`
	Name           string          `json:"name,omitempty"`
	DisplayName    string          `json:"display_name,omitempty"`
	Type           string          `json:"type,omitempty"`
	Prio           int             `json:"prio,omitempty"`
	Content        string          `json:"content,omitempty"`
	DisplayContent string          `json:"display_content,omitempty"`
	TTL            string          `json:"ttl,omitempty"`
	TTLCeil        int             `json:"ttl_ceil,omitempty"`
	SSLID          int             `json:"ssl_id,omitempty"`
	SSLStatus      int             `json:"ssl_status,omitempty"`
	SSLExpiresOn   int             `json:"ssl_expires_on,omitempty"`
	AutoTTL        int             `json:"auto_ttl,omitempty"`
	ServiceMode    string          `json:"service_mode,omitempty"`
	Properties     *RecordProperty `json:"props,omitempty"`
}

type RecordProperty struct {
	Proxiable   int `json:"proxiable,omitempty"`
	CloudOn     int `json:"cloud_on,omitempty"`
	CFOpen      int `json:"cf_open,omitempty"`
	SSL         int `json:"ssl,omitempty"`
	ExpiredSSL  int `json:"expired_ssl,omitempty"`
	ExpiringSSL int `json:"expiring_ssl,omitempty"`
	PendingSSL  int `json:"pending_ssl,omitempty"`
}

func (c *ClientAPI) RecLoadAll(z string) ([]Record, error) {
	data := &url.Values{
		"tkn":   {c.client.config.Token},
		"email": {c.client.config.Email},
		"a":     {"rec_load_all"},
		"z":     {z},
	}

	req, err := c.NewRequest(data)
	if err != nil {
		return nil, err
	}

	v := new(RecLoadAllResponse)
	_, err = c.client.Do(req, v)
	if err != nil {
		return nil, err
	}
	if v.Result == "error" || v.ErrorCode != "" {
		return nil, fmt.Errorf("%v %v: %v", req.Method, req.URL, v.Message)
	}

	return v.Response.Records.Objects, nil
}
