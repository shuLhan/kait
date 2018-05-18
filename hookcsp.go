// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type hookCSP struct {
	prefix string
	mmChan string
	ch     chan *CSPReport
}

//
// newHookCSP create, initialize, and return CSP hook that can be used as HTTP
// handler.
//
func newHookCSP() *hookCSP {
	return &hookCSP{
		ch: make(chan *CSPReport, 64),
	}
}

//
// ServeHTTP consume CSP report from HTTP.
//
func (hcsp *hookCSP) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("hookCSP.ServeHTTP:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cspWrapper := new(CSPReportWrapper)

	err = json.Unmarshal(reqBody, cspWrapper)
	if err != nil {
		log.Println("hookCSP.ServeHTTP:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	hcsp.ch <- &cspWrapper.CSPReport

	res.WriteHeader(http.StatusNoContent)
}
