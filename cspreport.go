// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

import (
	"bytes"
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
// CSPReportWrapper define the CSP wrapper for POST body.
//
type CSPReportWrapper struct {
	CSPReport CSPReport `json:"csp-report"`
}

//
// MarshalKV convert the report into `key=value key2=value2 ...` format.
//
func (report *CSPReport) MarshalKV() string {
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

	return buf.String()
}
