package main

import (
	word_filter "github.com/cyuxlif/go-word-filter"
	"flag"
)

var (
	configPath  = flag.String("c", "", `config path`)
)

func main(){
	options := word_filter.NewOptions(*configPath)
	filter := word_filter.New(options)
	filter.Run(options.TCPAddr)
}
