// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

type messageFormat uint

const (
	msgFormatKV messageFormat = 1 << iota
	msgFormatJSON
)

type messageMode uint

const (
	msgModeCSP messageMode = 1 << iota
)

type message struct {
	mode    messageMode
	content messageContent
}

func (msg *message) getContent(format messageFormat) string {
	var (
		bb  []byte
		err error
	)
	switch format {
	case msgFormatKV:
		bb, err = msg.content.MarshalKV()
	case msgFormatJSON:
		bb, err = msg.content.MarshalJSON()
	}

	if err != nil {
		return ""
	}

	return string(bb)
}
