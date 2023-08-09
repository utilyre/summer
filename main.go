package main

import (
	"log"
	"os"

	"gihtub.com/utilyre/summer/commands"
)

func main() {
	log.SetFlags(0)
	if err := commands.Execute(os.Args[1:]); err != nil {
		log.Fatalf("summer: %s", err)
	}
}
