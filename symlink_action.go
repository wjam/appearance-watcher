package main

import "os"

type symLinkFile struct {
	Path  string
	Light string
	Dark  string
}

func (s symLinkFile) ToLight() error {
	return s.createSymlinkAtomically(s.Light)
}

func (s symLinkFile) ToDark() error {
	return s.createSymlinkAtomically(s.Dark)
}

// createSymlinkAtomically will update the target symlink without deleting it first, so the symlink will always exist.
func (s symLinkFile) createSymlinkAtomically(dest string) error {
	tmp := s.Path + ".tmp"
	if err := os.Remove(tmp); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Symlink(dest, tmp); err != nil {
		return err
	}
	return os.Rename(tmp, s.Path)
}
