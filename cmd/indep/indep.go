package main

import (
	"os"

	"github.com/b-harvest/indep_node_alarm_go/cmd/indep/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
