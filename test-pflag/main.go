package main

import (
	"fmt"
	"github.com/spf13/pflag"
)

var cfgFile string

func main() {
	pflag.StringVarP(&cfgFile, "config", "c", "config.yml", "config file name")
	pflag.Parse()
	fmt.Println("config file name:", cfgFile)
}
