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
	msgMode messageMode
	webhook string
	channel string
	cl      *http.Client
	in      chan *message
	running chan bool
}

//
// NewMattermost create, initialize, and return forwarder for Mattermost.
//
func NewMattermost(webhook, channel string) Forwarder {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gas
		},
	}

	mm := &mattermost{
		msgMode: msgModeKeyValue,
		webhook: webhook,
		channel: channel,
		cl: &http.Client{
			Timeout:   defTimeout,
			Transport: tr,
		},
		in:      make(chan *message, 64),
		running: make(chan bool),
	}

	return mm
}

//
// Forward message into mattermost
//
func (mm *mattermost) Forward(msg *message) {
	mm.in <- msg
}

//
// MessageMode return the desired message mode to be forwarded.
//
func (mm *mattermost) MessageMode() messageMode {
	return mm.msgMode
}

//
// Start the mattermost forwarder.
//
func (mm *mattermost) Start() {
	running := true

	for running {
		select {
		case msg := <-mm.in:
			mm.forward(msg)
		case running = <-mm.running:
		}
	}
}

//
// Stop the mattermost forwarder.
//
func (mm *mattermost) Stop() {
	mm.running <- false
}

func (mm *mattermost) forward(msg *message) {
	if msg == nil {
		return
	}
	if len(msg.content) == 0 {
		return
	}

	content := text.StringJSONEscape(msg.content)

	var buf bytes.Buffer

	buf.WriteByte('{')
	buf.WriteString(`"channel":"`)
	if len(msg.channel) > 0 {
		buf.WriteString(msg.channel)
	} else {
		buf.WriteString(mm.channel)
	}

	buf.WriteString(`","text":"`)
	buf.WriteString(content)
	buf.WriteString(`"}`)

	req, err := http.NewRequest(http.MethodPost, mm.webhook, &buf)
	if err != nil {
		log.Println("ForwardCSP: http.NewRequest:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := mm.cl.Do(req)
	if err != nil {
		log.Println("ForwardCSP: cl.Do:", err)
		return
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ForwardCSP: ioutil.ReadAll:", err)
	}

	err = res.Body.Close()
	if err != nil {
		log.Println("ForwardCSP: res.Body.Close:", err)
	}
}
