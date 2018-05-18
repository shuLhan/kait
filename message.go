// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

type messageMode uint

const (
	msgModeKeyValue messageMode = 1 << iota
	msgModeJSON
)

type message struct {
	mode    messageMode
	channel string
	content string
}
