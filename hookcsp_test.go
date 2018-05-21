// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shuLhan/share/lib/ini"
	"github.com/shuLhan/share/lib/test"
)

func testHookCSPServeHTTP(t *testing.T) {
	cfg := &ini.Section{
		Vars: []*ini.Variable{{
			Key:      "mattermost-endpoint",
			KeyLower: "mattermost-endpoint",
			Value:    "http://172.31.3.185:8065/hooks/t934mqsdm3rxpytuhzfearyifc",
		}, {
			Key:      "mattermost-channel",
			KeyLower: "mattermost-channel",
			Value:    "log_test",
		}},
	}
	hookCSPTest := newHookCSP(cfg)

	cases := []struct {
		desc    string
		reqBody string
		exp     *message
	}{{
		desc: "Valid CSP body",
		reqBody: `{
  "csp-report": {
    "document-uri": "http://example.com/signup.html",
    "referrer": "",
    "blocked-uri": "http://example.com/css/style.css",
    "violated-directive": "style-src cdn.example.com",
    "original-policy": "default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports"
  }
}
`,
		exp: &message{
			mode:    msgModeCSP,
			content: cspReportTest[0],
		},
	}}

	for _, c := range cases {
		t.Log(c.desc)

		body := strings.NewReader(c.reqBody)

		req := httptest.NewRequest(http.MethodPost, defCSPPrefix, body)

		res := httptest.NewRecorder()

		hookCSPTest.ServeHTTP(res, req)

		got := <-hookCSPTest.chMsg

		test.Assert(t, "", c.exp, got, true)
	}
}

func TestHookCSP(t *testing.T) {
	t.Run("ServeHTTP", testHookCSPServeHTTP)
}
