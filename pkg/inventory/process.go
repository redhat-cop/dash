package inventory

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	cp "github.com/redhat-cop/dash/pkg/copy"
)

func (i *Inventory) Process(ns *string) error {

	if i.Namespace != "" {
		ns = &i.Namespace
	}

	if i.ResourceGroups != nil {
		for _, rg := range i.ResourceGroups {
			err := rg.ProcessResourceGroup(i.Prefix, ns)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (rg *ResourceGroup) ProcessResourceGroup(prefix string, ns *string) error {

	if rg.Namespace != "" {
		ns = &rg.Namespace
	}

	file, err := ioutil.TempDir("", "dash")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(file)

	if rg.Resources != nil {
		for _, r := range rg.Resources {
			err := r.ProcessResource(ns, file, prefix)
			if err != nil {
				return err
			}
		}
	}

	err = Reconcile(file, ns)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resource) ProcessResource(ns *string, f string, p string) error {

	if r.Namespace != "" {
		ns = &r.Namespace
	}

	log.Println("Resource: " + r.Name + ", Namespace: " + *ns)

	// TODO: Determine the type of resource
	if r.File != "" {
		err := r.ProcessFile(ns, f, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func Reconcile(path string, ns *string) error {

	p := path
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	cmdArgs := []string{"apply", "-f", filepath.Clean(abs)}
	if *ns != "" {
		cmdArgs = append(cmdArgs, "-n", *ns)
	}

	cmd := exec.Command("kubectl", cmdArgs...)
	fmt.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}
	fmt.Printf("%s\n", stdoutStderr)

	return nil

}

func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	log.Printf("Source: %s, Dest: %s\n", src, dst)
	if sourceFileStat.IsDir() {
		err = cp.CopyDir(src, dst)
		if err != nil {
			return err
		}
		return nil
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	err = cp.CopyFile(src, dst)
	if err != nil {
		return err
	}
	return nil
}
