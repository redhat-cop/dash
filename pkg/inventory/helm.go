package inventory

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// HelmChart manages pulling and processing of Helm charts
type HelmChart struct {
	Chart      string      `yaml:"chart"`
	URL        string      `yaml:"url"`
	Values     interface{} `yaml:"values"`
	ValueFiles []string    `yaml:"valueFiles"`
}

// Process runs `helm fetch` followed by `helm template`
func (h *HelmChart) Process(ns *string, r *Resource) error {

	// set values
	prefix, err := filepath.Abs(r.Prefix)
	if err != nil {
		return err
	}
	output, err := filepath.Abs(r.Output)
	if err != nil {
		return err
	}

	// validate chart name
	chart := h.getName()
	if chart == "" {
		return fmt.Errorf("chart validation failed. Here's the chart: %s", chart)
	}

	// fetch chart if from URL
	// helm fetch --untar --untardir . 'stable/redis'
	cmdArgs := []string{"fetch", "--untar", "--untardir", output + "/charts", h.Chart}
	cmd := exec.Command("helm", cmdArgs...)
	log.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}
	log.Printf("%s\n", stdoutStderr)

	// ensure output directory exists
	outputDir := output + "/" + string(r.Action)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	// helm template --output-dir './redis-final' './redis' --set
	// base arguments to `helm template`
	cmdArgs = []string{"template", "--output-dir", outputDir, output + "/charts/" + chart}
	// generate flags for valueFiles
	for _, f := range h.ValueFiles {
		cmdArgs = append(cmdArgs, "-f", prefix+"/"+f)
	}
	// write embedded values to file and pass as arg
	if h.Values != nil {
		vOut := output + "/charts/" + chart + "/dash_values.yaml"
		err = marshalValues(h.Values, vOut)
		if err != nil {
			return err
		}
		cmdArgs = append(cmdArgs, "-f", vOut)
	}

	// execute helm command and handle output
	cmd = exec.Command("helm", cmdArgs...)
	log.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}
	log.Printf("%s\n", stdoutStderr)

	return nil

}

func (h *HelmChart) getName() string {
	var canonical = regexp.MustCompile(`^[a-zA-Z]+/[a-zA-Z]+$`)
	var url = regexp.MustCompile(`^http[s]?:\/\/.*$`)

	switch {
	case canonical.MatchString(h.Chart):
		return strings.Split(h.Chart, "/")[1]
	case url.MatchString(h.Chart):
		c := strings.Split(h.Chart, "/")
		return c[len(c)-1]
	default:
		return ""
	}
}

func marshalValues(v interface{}, out string) error {
	y, err := yaml.Marshal(&v)
	if err != nil {
		return err
	}

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(y))

	return err
}
