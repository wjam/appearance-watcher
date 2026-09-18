package main

import "os"

type writeToFile struct {
	Path  string
	Light string
	Dark  string
}

func (w writeToFile) ToLight() error {
	return os.WriteFile(w.Path, []byte(w.Light), 0600)
}

func (w writeToFile) ToDark() error {
	return os.WriteFile(w.Path, []byte(w.Dark), 0600)
}
