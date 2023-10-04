package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"tools/clean/atlas"
	"tools/clean/pe"
	"tools/clean/vpc"
)

func main() {
	if err := Clean(os.Args); err != nil {
		log.Printf("Invocation failed: %s", err)
		log.Fatalf("Usage: %s {atlas|pe|vpc}", os.Args[0])
	}
}

var cleanAtlas = atlas.CleanAtlas

var cleanPEs = pe.CleanPEs

var cleanVPCs = vpc.CleanVPCs

func Clean(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Wrong number of arguments: expected 1 got %d", len(args)-1)
	}
	action := strings.ToLower(args[1])
	switch action {
	case "atlas":
		cleanAtlas()
	case "pe":
		cleanPEs()
	case "vpc":
		cleanVPCs()
	default:
		return fmt.Errorf("Unsupported action %q", action)
	}
	return nil
}
