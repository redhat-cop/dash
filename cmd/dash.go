package main

import (
	"flag"
	"fmt"

	inv "github.com/redhat-cop/dash/pkg/inventory"
)

var invPath string

func init() {
	flag.StringVar(&invPath, "i", "./", "Path to Inventory (relative or absolute); Defaults to ./")
	flag.Parse()
}

func main() {

	var i inv.Inventory
	var ns string
	i.Load(invPath)
	err := i.Process(&ns)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

}
