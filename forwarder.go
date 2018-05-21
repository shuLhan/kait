// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kait

//
// Forwarder is an interface that implement message forwarder from Kait to
// other, external, consumer.
//
type Forwarder interface {
	Forward(*message, func())
}
