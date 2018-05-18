// Copyright 2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Kait is generic HTTP proxy that forwards hooks from external service to
// internal service.
//
package main

import (
	"flag"
	"log"

	"github.com/shuLhan/kait"
)

const (
	defConfigFile = "/etc/kait.conf"
)

func main() {
	var cfgFile = defConfigFile

	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		cfgFile = args[0]
	}

	server := &kait.Kait{}

	err := server.Init(cfgFile)
	if err != nil {
		log.Fatal("Init:" + err.Error())
	}

	err = server.Start()
	if err != nil {
		log.Print(err)
	}
}
