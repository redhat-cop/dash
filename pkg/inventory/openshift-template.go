package inventory

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type OpenShiftTemplate struct {
	Template   string            `yaml:"template"`
	Params     map[string]string `yaml:"params"`
	ParamFiles []string          `yaml:"paramFiles"`
}

func (ot *OpenShiftTemplate) Process(ns *string, r *Resource) error {

	p := r.Prefix + "/" + ot.Template
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	// oc process -f template-file -p PARAM=foo --param-file
	cmdArgs := []string{"process", "--local", "-f", abs}
	for key, param := range ot.Params {
		cmdArgs = append(cmdArgs, "-p", key+"="+param)
	}
	for _, pf := range ot.ParamFiles {
		cmdArgs = append(cmdArgs, "--param-file", pf)
	}
	cmd := exec.Command("oc", cmdArgs...)
	log.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}

	// write resulting resource to file
	output_dir := filepath.Clean(r.Output + "/" + string(r.Action))
	out, err := os.Create(output_dir + "/" + ot.Template)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	log.Printf("wrote %s\n", out.Name())
	_, err = out.Write(stdoutStderr)
	if err != nil {
		return err
	}

	return nil
}
