package main

import (
	"log"
	"os"

	"pharos-exporter/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
