package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	inv "github.com/redhat-cop/dash/pkg/inventory"
)

const (
	invPathDefault = "./"
	invPathUsage   = "Path to Inventory, relative or absolute"
)

var (
	version     string
	invPath     string
	showVersion bool
)

func init() {
	flag.StringVar(&invPath, "inventory", invPathDefault, invPathUsage)
	flag.StringVar(&invPath, "i", invPathDefault, invPathUsage+" (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "See version")
	flag.Parse()
}

func main() {

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	var i inv.Inventory
	var ns string

	yamlFile, err := ioutil.ReadFile(invPath + "dash.yaml")
	if err != nil {
		fmt.Printf("Error: Couldn't load dash inventory: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	i.Load(yamlFile, invPath)
	err = i.Process(&ns)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		flag.Usage()
		os.Exit(1)
	}

}
