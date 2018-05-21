// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"bytes"
	"encoding/json"
)

//
// CSPReport define a Content-Security-Report.
//
type CSPReport struct {
	DocumentURI       string `json:"document-uri"`
	Referrer          string `json:"referrer"`
	BlockedURI        string `json:"blocked-uri"`
	ViolatedDirective string `json:"violated-directive"`
	OriginalPolicy    string `json:"original-policy"`
}

//
// CSPWrapper define the CSP wrapper for POST body.
//
type CSPWrapper struct {
	Report *CSPReport `json:"csp-report"`
}

//
// MarshalKV convert the report into `key=value key2=value2 ...` format.
//
func (report *CSPReport) MarshalKV() ([]byte, error) {
	var buf bytes.Buffer

	_, _ = buf.WriteString("document-uri=")
	_, _ = buf.WriteString(report.DocumentURI)
	_, _ = buf.WriteString(" referrer=")
	_, _ = buf.WriteString(report.Referrer)
	_, _ = buf.WriteString(" blocked-uri=")
	_, _ = buf.WriteString(report.BlockedURI)
	_, _ = buf.WriteString(" document-uri=")
	_, _ = buf.WriteString(report.DocumentURI)
	_, _ = buf.WriteString(" violated-directive=")
	_, _ = buf.WriteString(report.ViolatedDirective)
	_, _ = buf.WriteString(" original-policy=")
	_, _ = buf.WriteString(report.OriginalPolicy)

	return buf.Bytes(), nil
}

//
// MarshalJSON convert the report into JSON.
//
func (report *CSPReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(report)
}
