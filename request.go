package main

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var client = http.Client{Transport: transport()}

func getFile(rawUrl, name string) error {
	oldFile := dir + name
	if isExist(oldFile) {
		if err := os.Remove(oldFile); err != nil {
			return err
		}
	}
	fi, err := os.Create(name)
	if err != nil {
		return err
	}
	resp, err := client.Get(rawUrl)
	if err != nil {
		return err
	}
	_, err = io.Copy(fi, resp.Body)
	return resp.Body.Close()
}

func transport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		DialContext: (&net.Dialer{
			Timeout: 180 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 60 * time.Second,
		ForceAttemptHTTP2:   true,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 50,
	}
}
