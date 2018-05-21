// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/shuLhan/share/lib/text"
)

const (
	defTimeout = 5 * time.Second
)

type mattermost struct {
	msgFormat messageFormat
	endpoint  string
	channel   string
	cl        *http.Client
}

//
// NewMattermost create, initialize, and return forwarder for Mattermost.
//
func NewMattermost(endpoint, channel string) Forwarder {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gas
		},
	}

	mm := &mattermost{
		msgFormat: msgFormatKV,
		endpoint:  endpoint,
		channel:   channel,
		cl: &http.Client{
			Timeout:   defTimeout,
			Transport: tr,
		},
	}

	return mm
}

//
// Forward message into mattermost
//
func (mm *mattermost) Forward(msg *message, cbOnFinish func()) {
	var (
		err     error
		req     *http.Request
		res     *http.Response
		content string
		buf     bytes.Buffer
	)

	if msg == nil {
		goto out
	}

	content = msg.getContent(mm.msgFormat)
	if len(content) == 0 {
		goto out
	}

	content = text.StringJSONEscape(content)

	buf.WriteByte('{')
	if len(mm.channel) > 0 {
		buf.WriteString(`"channel":"`)
		buf.WriteString(mm.channel)
		buf.WriteString(`",`)
	}

	buf.WriteString(`"text":"`)
	buf.WriteString(content)
	buf.WriteString(`"}`)

	req, err = http.NewRequest(http.MethodPost, mm.endpoint, &buf)
	if err != nil {
		log.Println("ForwardCSP: http.NewRequest:", err)
		goto out
	}

	req.Header.Set("Content-Type", "application/json")

	res, err = mm.cl.Do(req)
	if err != nil {
		log.Println("ForwardCSP: cl.Do:", err)
		goto out
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ForwardCSP: ioutil.ReadAll:", err)
	}

	err = res.Body.Close()
	if err != nil {
		log.Println("ForwardCSP: res.Body.Close:", err)
	}

out:
	cbOnFinish()
}
