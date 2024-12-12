package main

import (
	"log"
	"os"

	"github.com/koron-go/subcmd"
	"github.com/koron/giftool/internal/info"
)

var rootSet = subcmd.DefineRootSet(
	info.Info,
)

func main() {
	if err := subcmd.Run(rootSet, os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}