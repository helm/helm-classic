package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Args are: %v\n", os.Args)
	fmt.Printf("Helm home is: %s\n", os.Getenv("HELM_HOME"))
	fmt.Printf("Helm command is: %s\n", os.Getenv("HELM_COMMAND"))
	fmt.Printf("Helm default repo is: %s\n", os.Getenv("HELM_DEFAULT_REPO"))
}
