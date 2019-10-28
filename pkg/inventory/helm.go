package inventory

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type HelmChart struct {
	Chart      string      `yaml:"chart"`
	Url        string      `yaml:"url"`
	Values     interface{} `yaml:"values"`
	ValueFiles []string    `yaml:"valueFiles"`
}

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
		return errors.New(fmt.Sprintf("chart validation failed. Here's the chart: %s", chart))
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
	output_dir := output + "/" + string(r.Action)
	if _, err := os.Stat(output_dir); os.IsNotExist(err) {
		os.Mkdir(output_dir, os.ModePerm)
	}

	// helm template --output-dir './redis-final' './redis' --set
	// base arguments to `helm template`
	cmdArgs = []string{"template", "--output-dir", output_dir, output + "/charts/" + chart}
	// generate flags for valueFiles
	for _, f := range h.ValueFiles {
		cmdArgs = append(cmdArgs, "-f", prefix+"/"+f)
	}
	// write embedded values to file and pass as arg
	if h.Values != nil {
		v_out := output + "/charts/" + chart + "/dash_values.yaml"
		err = marshalValues(h.Values, v_out)
		if err != nil {
			return err
		}
		cmdArgs = append(cmdArgs, "-f", v_out)
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
