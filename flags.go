package main

import (
	"flag"
	"time"

	"github.com/ostrost/ostent/types"
)

// periodFlag is a minimal refresh period
var periodFlag = types.PeriodValue{Duration: types.Duration(time.Second)} // default

// TODO This ought to go into a command
// collectdBindFlag is a bindValue holding the ostent collectd bind address.
// var collectdBindFlag = types.NewBindValue("", "8051") // "" by default meaning DO NOT BIND

func init() {
	flag.Var(&periodFlag, "u", "Collection (update) interval")
	flag.Var(&periodFlag, "update", "Collection (update) interval")
	// flag.Var(&collectdBindFlag, "collectdb",    "short for collectdbind")
	// flag.Var(&collectdBindFlag, "collectdbind", "Bind address for collectd receiving")
}
