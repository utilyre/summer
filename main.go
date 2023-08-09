package main

import (
	"log"
	"os"

	"gihtub.com/utilyre/summer/command"
)

func main() {
	log.SetFlags(0)
	if err := command.Execute(os.Args[1:]); err != nil {
		log.Fatalf("summer: %s", err)
	}
}
