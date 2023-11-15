package main

import (
	"github.com/diogosilva96/etf-cli/cmd"
	"log"
)

func main() {
	cli, err := cmd.NewCli()
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}
