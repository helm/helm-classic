package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Args are: %v\n", os.Args)
	// Although helmc itself may use the new HELMC_HOME environment variable to optionally define its
	// home directory, to maintain compatibility with charts created for the ORIGINAL helm, helmc
	// still expands "legacy" Helm environment variables, which Helm Classic plugins continue to use.
	fmt.Printf("Helm home is: %s\n", os.Getenv("HELM_HOME"))
	fmt.Printf("Helm command is: %s\n", os.Getenv("HELM_COMMAND"))
	fmt.Printf("Helm default repo is: %s\n", os.Getenv("HELM_DEFAULT_REPO"))
}
