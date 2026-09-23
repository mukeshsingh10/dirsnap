package main

import "path/filepath"

type File struct {
	Name     string
	Path     string
	Metadata *FileMD
}

func NewFile(path string) *File {
	if path == "" {
		path = "."
	}
	filename := filepath.Base(path)
	metadata := NewFileMD(filename)

	return &File{
		Name:     filename,
		Path:     path,
		Metadata: metadata,
	}
}
