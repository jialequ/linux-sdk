package main

import (
	"github.com/jialequ/linux-sdk/core/load"
	"github.com/jialequ/linux-sdk/core/logx"
	"github.com/jialequ/linux-sdk/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
