package inventory

import (
	"path/filepath"
)

// FileTemplate static manifest file type
type FileTemplate struct {
	Path string `yaml:"path"`
}

// Process moves static file resources to Action directory
func (ft *FileTemplate) Process(ns *string, r *Resource) error {

	p := r.Prefix + "/" + ft.Path
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	// Copy file to tmpdir
	err = copy(abs, r.Output+"/"+string(r.Action))
	if err != nil {
		return err
	}

	return nil
}
