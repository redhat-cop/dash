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
)

type HelmChart struct {
	Chart string                 `yaml:"chart"`
	Url   string                 `yaml:"url"`
	Set   map[string]interface{} `yaml:"set"`
}

func (h *HelmChart) Process(ns *string, r *Resource) error {

	// validate chart name
	chart := h.getName()
	if chart == "" {
		return errors.New(fmt.Sprintf("chart validation failed. Here's the chart: %s", chart))
	}

	// fetch chart if from URL
	// helm fetch --untar --untardir . 'stable/redis'
	cmdArgs := []string{"fetch", "--untar", "--untardir", filepath.Clean(r.Output + "/charts"), h.Chart}
	cmd := exec.Command("helm", cmdArgs...)
	log.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}
	log.Printf("%s\n", stdoutStderr)

	// ensure output directory exists
	output_dir := filepath.Clean(r.Output + "/" + string(r.Action))
	if _, err := os.Stat(output_dir); os.IsNotExist(err) {
		os.Mkdir(output_dir, os.ModePerm)
	}

	// helm template --output-dir './redis-final' './redis' --set
	cmdArgs = []string{"template", "--output-dir", output_dir, filepath.Clean(r.Output + "/charts/" + chart)}
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
	var canonical = regexp.MustCompile(`[a-zA-Z]+/[a-zA-Z]+`)

	switch {
	case canonical.MatchString(h.Chart):
		return strings.Split(h.Chart, "/")[1]
	default:
		return ""
	}
}
