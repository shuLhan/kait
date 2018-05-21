// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"testing"

	"github.com/shuLhan/share/lib/test"
)

var cspReportTest = []*CSPReport{{
	DocumentURI:       `http://example.com/signup.html`,
	Referrer:          ``,
	BlockedURI:        `http://example.com/css/style.css`,
	ViolatedDirective: `style-src cdn.example.com`,
	OriginalPolicy:    `default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports`,
}, {
	DocumentURI:       `http://example.com/signup.html`,
	Referrer:          `te"st`,
	BlockedURI:        `http://example.com/css/style.css`,
	ViolatedDirective: `style-src cdn.example.com`,
	OriginalPolicy:    `default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports`,
}}

func TestMarshalKV(t *testing.T) {
	cases := []struct {
		desc string
		in   *CSPReport
		exp  string
	}{{
		desc: "With normal value",
		in:   cspReportTest[0],
		exp:  `document-uri=http://example.com/signup.html referrer= blocked-uri=http://example.com/css/style.css document-uri=http://example.com/signup.html violated-directive=style-src cdn.example.com original-policy=default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports`,
	}, {
		desc: "With escaped characters",
		in:   cspReportTest[1],
		exp:  `document-uri=http://example.com/signup.html referrer=te"st blocked-uri=http://example.com/css/style.css document-uri=http://example.com/signup.html violated-directive=style-src cdn.example.com original-policy=default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports`,
	}}

	for _, c := range cases {
		t.Log(c.desc)

		got, err := c.in.MarshalKV()
		if err != nil {
			t.Fatal(err)
		}

		test.Assert(t, "cspreport", c.exp, string(got), true)
	}
}
