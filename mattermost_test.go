// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"testing"

	"github.com/shuLhan/share/lib/ini"
)

const (
	testConfig = "test.conf"
)

var (
	mm Forwarder
)

func testForwardCSP(t *testing.T) {
	cases := []struct {
		desc    string
		channel string
		in      *CSPReport
	}{{
		desc: "With normal CSP report",
		in: &CSPReport{
			DocumentURI:       `http://example.com/signup.html`,
			Referrer:          `test`,
			BlockedURI:        `http://example.com/css/style.css`,
			ViolatedDirective: `style-src cdn.example.com`,
			OriginalPolicy:    `default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports`,
		},
	}, {
		desc:    "With escaped char",
		channel: "log_virtuals",
		in: &CSPReport{
			DocumentURI:       `http://example.com/signup.html/?key=val%20ue`,
			Referrer:          `test`,
			BlockedURI:        `http://example.com/css/style.css`,
			ViolatedDirective: `style-src cdn.example.com`,
			OriginalPolicy:    `default-src "none"; style-src cdn.example.com; report-uri /_/csp-reports`,
		},
	}}

	for _, c := range cases {
		t.Log(c.desc)

		msg := &message{
			mode:    mm.MessageMode(),
			channel: c.channel,
			content: c.in.MarshalKV(),
		}

		mm.Forward(msg)
	}
}

func TestMattermost(t *testing.T) {
	env, err := ini.Open(testConfig)
	if err != nil {
		return
	}

	webhook, _ := env.Get(sectionFW, subMM, keyWebhook)
	channel, _ := env.Get(sectionFW, subMM, keyDefChan)

	mm = NewMattermost(webhook, channel)

	go mm.Start()

	t.Run("ForwardCSP", testForwardCSP)

	mm.Stop()
}
