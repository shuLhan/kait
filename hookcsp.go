// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/shuLhan/share/lib/ini"
)

type hookCSP struct {
	prefix  string
	chMsg   chan *message
	running chan bool
	fwds    []Forwarder
}

//
// newHookCSP create, initialize, and return CSP hook that can be used as HTTP
// handler.
//
func newHookCSP(cfg *ini.Section) (hook *hookCSP) {
	hook = &hookCSP{
		chMsg:   make(chan *message, maxMessageBuffer),
		running: make(chan bool, 1),
	}

	hook.prefix, _ = cfg.Get(keyPrefix, defCSPPrefix)
	mmChan, _ := cfg.Get(keyMMChan, "")
	mmEndpoint, _ := cfg.Get(keyMMEndpoint, "")

	if len(mmEndpoint) > 0 {
		fw := NewMattermost(mmEndpoint, mmChan)
		hook.fwds = append(hook.fwds, fw)
	}

	if len(hook.fwds) == 0 {
		hook = nil
		return
	}

	return
}

//
// Start forwarding message.
//
func (hook *hookCSP) Start() {
	var msg *message
	running := true
	for running {
		select {
		case msg = <-hook.chMsg:
			cur := atomic.LoadInt32(&nRoutines)
			if cur >= maxRoutines {
				continue
			}

			for _, fw := range hook.fwds {
				atomic.AddInt32(&nRoutines, 1)

				go fw.Forward(msg, func() {
					atomic.AddInt32(&nRoutines, -1)
				})
			}
		case running = <-hook.running:
		}
	}
}

//
// Stop forwarding message.
//
func (hook *hookCSP) Stop() {
	hook.running <- false
}

//
// ServeHTTP consume CSP report from HTTP.
//
func (hook *hookCSP) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("hookCSP.ServeHTTP:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cspWrapper := new(CSPWrapper)

	err = json.Unmarshal(reqBody, cspWrapper)
	if err != nil {
		log.Println("hookCSP.ServeHTTP:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("%+v\n", cspWrapper.Report)

	msg := &message{
		mode:    msgModeCSP,
		content: cspWrapper.Report,
	}

	hook.chMsg <- msg

	cspWrapper.Report = nil

	res.WriteHeader(http.StatusNoContent)
}
