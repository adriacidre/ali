package main

import (
	"log"
	"os"

	"github.com/adriacidre/aliases/config"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("homedir: %s", err)
	}

	path := home + "/.aliases"
	c := config.Config{}
	c.Parse(path)
	// c.Save(path)

	for _, c := range c.Commands {
		println("  " + c.Usage + ": " + c.Description)
	}
	println("")
}
