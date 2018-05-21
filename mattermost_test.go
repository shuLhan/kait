// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"sync"
	"testing"

	"github.com/shuLhan/share/lib/ini"
)

const (
	testConfig = "test.conf"
)

var (
	mm     Forwarder
	testWG sync.WaitGroup
)

func testForwardCSP(t *testing.T) {
	cases := []struct {
		desc string
		msg  *message
	}{{
		desc: "With normal CSP report",
		msg: &message{
			mode: msgModeCSP,
			content: &CSPReport{
				DocumentURI:       `http://example.com/signup.html`,
				Referrer:          `test`,
				BlockedURI:        `http://example.com/css/style.css`,
				ViolatedDirective: `style-src cdn.example.com`,
				OriginalPolicy:    `default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports`,
			},
		},
	}, {
		desc: "With escaped char",
		msg: &message{
			mode: msgModeCSP,
			content: &CSPReport{
				DocumentURI:       `http://example.com/signup.html/?key=val%20ue`,
				Referrer:          `test`,
				BlockedURI:        `http://example.com/css/style.css`,
				ViolatedDirective: `style-src cdn.example.com`,
				OriginalPolicy:    `default-src "none"; style-src cdn.example.com; report-uri /_/csp-reports`,
			},
		},
	}}

	for _, c := range cases {
		t.Log(c.desc)

		testWG.Add(1)

		mm.Forward(c.msg, func() {
			testWG.Done()
		})
	}
}

func TestMattermost(t *testing.T) {
	env, err := ini.Open(testConfig)
	if err != nil {
		return
	}

	endpoint, _ := env.Get(sectionCSP, "", keyMMEndpoint)
	channel, _ := env.Get(sectionCSP, "", keyMMChan)

	mm = NewMattermost(endpoint, channel)

	t.Run("ForwardCSP", testForwardCSP)

	testWG.Wait()
}
