package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

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

	yamlFile, err := ioutil.ReadFile(invPath + "dash.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	i.Load(yamlFile, invPath)
	err = i.Process(&ns)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

}
