package main

import (
	"github.com/OpenBankingUK/conformance-suite/pkg/discovery/templates"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Generate a draft generic discovery")
	err := templates.DraftDiscovery()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}
}
