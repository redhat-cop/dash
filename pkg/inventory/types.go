package inventory

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const ACTION string = "apply"

type DashMeta struct {
	Prefix string `yaml:"prefix"`
	Output string `yaml:"output"`
	Action string `yaml:"action"`
}

type Inventory struct {
	DashMeta       `yaml:",inline"`
	Version        int64           `yaml:"version"`
	Namespace      string          `yaml:"namespace"`
	ResourceGroups []ResourceGroup `yaml:"resource_groups"`
	Args           []string
}

type ResourceGroup struct {
	DashMeta  `yaml:",inline"`
	Name      string     `yaml:"name"`
	Namespace string     `yaml:"namespace"`
	Resources []Resource `yaml:"resources"`
}

type Resource struct {
	DashMeta          `yaml:",inline"`
	Name              string            `yaml:"name"`
	Namespace         string            `yaml:"namespace"`
	File              FileTemplate      `yaml:"file"`
	Helm              HelmChart         `yaml:"helm"`
	OpenShiftTemplate OpenShiftTemplate `yaml:"openshiftTemplate"`
}

type Template interface {
	Process(ns *string, r *Resource) error
}

func (i *Inventory) Load(yf []byte, pre string) *Inventory {

	file, err := ioutil.TempDir("", "dash")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(file)

	i.Prefix = pre
	i.Output = file

	err = yaml.Unmarshal(yf, &i)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	err = i.setDefaults()
	if err != nil {
		log.Fatalf("Failed setting defaults: %v", err)
	}

	log.Print(i)

	return i
}

func (i *Inventory) setDefaults() error {
	i.Action = ACTION
	for rgi, rg := range i.ResourceGroups {
		if rg.Action == "" {
			i.ResourceGroups[rgi].Action = ACTION
		}
		i.ResourceGroups[rgi].Prefix = i.Prefix
		i.ResourceGroups[rgi].Output = i.Output
		for ri, r := range rg.Resources {
			if r.Action == "" {
				i.ResourceGroups[rgi].Resources[ri].Action = ACTION
			}
			i.ResourceGroups[rgi].Resources[ri].Prefix = i.Prefix
			i.ResourceGroups[rgi].Resources[ri].Output = i.Output
		}
	}

	log.Print(i)
	return nil
}
