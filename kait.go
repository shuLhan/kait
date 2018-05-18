// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Package kait provide library for building generic HTTP hooks for forwarding
// external service to internal service.
//
package kait

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/shuLhan/share/lib/ini"
)

const (
	defListen    = "0.0.0.0:8417"
	defPrefix    = "/hooks"
	defCSPPrefix = "/csp"

	// List of config keys
	keyListen = "listen"
	keyPrefix = "prefix"

	sectionKait = "kait"
	sectionFW   = "forwarder"
	sectionCSP  = "csp"

	subMM = "mattermost"

	keyDefChan = "default-channel"
	keyWebhook = "webhook"
	keyMMChan  = "mattermost-channel"
)

var (
	defServerRTO            = 3 * time.Second
	defServerWTO            = 5 * time.Second
	defServerMaxHeaderBytes = 1024 * 1024 // 1 MB
)

//
// Kait is HTTP server that handle hooks response.
//
type Kait struct {
	prefix     string
	listenAddr string
	listener   net.Listener
	mux        *http.ServeMux
	server     *http.Server
	env        *ini.Ini
	hookCSP    *hookCSP
	fwds       []Forwarder
	running    chan bool
}

//
// Init will load Kait environment from file and create all enabled
// forwarders.
//
func (kait *Kait) Init(envFile string) (err error) {
	kait.env, err = ini.Open(envFile)
	if err != nil {
		return
	}

	kait.initServer()
	kait.initCSP()
	kait.initForwarders()

	return
}

func (kait *Kait) initServer() {
	sec := kait.env.GetSection(sectionKait, "")
	if sec == nil {
		kait.prefix = defPrefix
		kait.listenAddr = defListen
	} else {
		kait.prefix, _ = sec.Get(keyPrefix, defPrefix)
		kait.listenAddr, _ = sec.Get(keyListen, defListen)
	}

	kait.mux = http.NewServeMux()
	kait.server = &http.Server{
		Handler:        kait.mux,
		ReadTimeout:    defServerRTO,
		WriteTimeout:   defServerWTO,
		MaxHeaderBytes: defServerMaxHeaderBytes,
	}
}

func (kait *Kait) initCSP() {
	sec := kait.env.GetSection(sectionCSP, "")
	if sec == nil {
		return
	}

	kait.hookCSP = newHookCSP()

	kait.hookCSP.prefix, _ = sec.Get(keyPrefix, defCSPPrefix)
	kait.hookCSP.mmChan, _ = sec.Get(keyMMChan, "")

	kait.mux.Handle(kait.prefix+kait.hookCSP.prefix, kait.hookCSP)
}

func (kait *Kait) initForwarders() {
	kait.initMattermost()
}

func (kait *Kait) initMattermost() {
	sec := kait.env.GetSection(sectionFW, subMM)
	if sec == nil {
		return
	}

	webhook, _ := sec.Get(keyWebhook, "")
	defchan, _ := sec.Get(keyDefChan, "")

	if len(webhook) > 0 {
		fwMM := NewMattermost(webhook, defchan)
		kait.fwds = append(kait.fwds, fwMM)
	}
}

//
// Start will run all forwarders as go routine and start the HTTP server.
//
func (kait *Kait) Start() (err error) {
	for x := range kait.fwds {
		kait.fwds[x].Start()
	}

	kait.listener, err = net.Listen("tcp", kait.listenAddr)
	if err != nil {
		return
	}

	go kait.forwarder()

	err = kait.server.Serve(kait.listener)
	if err != nil {
		log.Println("Start.Serve:", err)
	}

	for x := range kait.fwds {
		kait.fwds[x].Stop()
	}

	return
}

//
// Stop the server.
//
func (kait *Kait) Stop() {
	kait.running <- false
}

func (kait *Kait) forwarder() {
	running := true

	for running {
		select {
		case cspReport := <-kait.hookCSP.ch:
			kait.forwardCSP(cspReport)
		case running = <-kait.running:
		}
	}
}

func (kait *Kait) forwardCSP(report *CSPReport) {
	var err error
	var bb []byte
	msg := new(message)

	for _, fw := range kait.fwds {
		switch fw.MessageMode() {
		case msgModeKeyValue:
			msg.mode = msgModeKeyValue
			msg.content = report.MarshalKV()

		case msgModeJSON:
			msg.mode = msgModeJSON
			bb, err = json.Marshal(report)
			if err != nil {
				log.Println("forwardCSP: json.Marshal:", err)
				return
			}
			msg.content = string(bb)
		}

		fw.Forward(msg)
	}
}
