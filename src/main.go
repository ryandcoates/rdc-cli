package main

import (
	"log"

	"rdc/cmd/rdc"
)

func main() {
	if err := rdc.Execute(); err != nil {
		log.Fatal(err)
	}
}
