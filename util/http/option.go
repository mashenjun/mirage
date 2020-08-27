package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func TimeoutOption(timeout time.Duration) func(client *Client) {
	return func(c *Client) {
		c.Client.Timeout = timeout
	}
}

func MaxIdleConnOption(connCount int) func(client *Client) {
	return func(c *Client) {
		// todo: is there better way to update transport
		// get the old transport and replace by new one
		oldTransport := c.Transport
		oldTransportPointer, ok := oldTransport.(*http.Transport)
		if !ok {
			panic(fmt.Sprintf("transport not an *http.Transport"))
		}
		// create new transport
		newTransport := &http.Transport{
			Proxy:                 oldTransportPointer.Proxy,
			DialContext:           oldTransportPointer.DialContext,
			MaxIdleConns:          connCount,
			MaxIdleConnsPerHost:   connCount,
			IdleConnTimeout:       oldTransportPointer.IdleConnTimeout,
			TLSHandshakeTimeout:   oldTransportPointer.TLSHandshakeTimeout,
			ExpectContinueTimeout: oldTransportPointer.ExpectContinueTimeout,
			// todo: set tls config
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.Transport = newTransport
	}
}
