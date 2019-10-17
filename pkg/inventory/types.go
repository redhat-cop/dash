package inventory

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Inventory struct {
	Version        int64           `yaml:"version"`
	Namespace      string          `yaml:"namespace"`
	ResourceGroups []ResourceGroup `yaml:"resource_groups"`
	Prefix         string
}

type ResourceGroup struct {
	Name      string     `yaml:"name"`
	Namespace string     `yaml:"namespace"`
	Resources []Resource `yaml:"resources"`
}

type Resource struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	File      string `yaml:"file"`
	Action    Action `yaml:"action"`
}

type Action string

// implement the Unmarshaler interface on Action, so we can default it to "apply"
func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*a = Action("apply")
	} else {
		*a = Action(s)
	}
	return nil
}

func (i *Inventory) Load(pre string) *Inventory {

	i.Prefix = pre

	yamlFile, err := ioutil.ReadFile(pre + "dash.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, i)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return i
}
