package main

import (
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
    l := lex("", "+jsonschema:validation:default=3")
    l.run()
	return nil
}
