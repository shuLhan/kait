// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Package kait provide library for building generic HTTP hooks for forwarding
// external service to internal service.
//
package kait

import (
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
	sectionCSP  = "csp"

	keyMMChan     = "mattermost-channel"
	keyMMEndpoint = "mattermost-endpoint"

	maxMessageBuffer = 64
	maxRoutines      = maxMessageBuffer * 2
)

var (
	defServerRTO            = 3 * time.Second
	defServerWTO            = 5 * time.Second
	defServerMaxHeaderBytes = 1024 * 1024 // 1 MB

	nRoutines int32
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
	cfg := kait.env.GetSection(sectionCSP, "")
	if cfg == nil {
		return
	}

	kait.hookCSP = newHookCSP(cfg)
	if kait.hookCSP == nil {
		return
	}

	kait.mux.Handle(kait.prefix+kait.hookCSP.prefix, kait.hookCSP)

	go kait.hookCSP.Start()
}

//
// Start will run all forwarders as go routine and start the HTTP server.
//
func (kait *Kait) Start() (err error) {
	kait.listener, err = net.Listen("tcp", kait.listenAddr)
	if err != nil {
		return
	}

	err = kait.server.Serve(kait.listener)
	if err != nil {
		log.Println("Start.Serve:", err)
	}

	return
}

//
// Stop the server.
//
func (kait *Kait) Stop() {
	kait.running <- false
}
