package inventory

import "path/filepath"

func (r *Resource) ProcessFile(ns *string, f string, p string) error {

	p = p + r.File
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	// Copy file to tmpdir
	err = copy(abs, f+string(r.Action))
	if err != nil {
		return err
	}

	return nil
}
