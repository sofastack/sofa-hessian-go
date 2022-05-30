package main

import (
	"log"

	"github.com/sofastack/sofa-hessian-go/sofahessian/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
