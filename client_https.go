package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

type httpclient struct {
	Net        string
	Timeout    time.Duration
	HTTPClient *http.Client
}

func (c *httpclient) Exchange(msg *dns.Msg, upstream string) (*dns.Msg, int, error) {
	data, err := msg.Pack()
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", upstream, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", dnsMsgType)

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("http error %d", resp.StatusCode)
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	m := new(dns.Msg)
	if err = m.Unpack(data); err != nil {
		return nil, 0, err
	}

	return m, 0, nil
}
