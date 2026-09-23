package main

import (
	"os"
	"syscall"
	"time"
)

type FileMD struct {
	Size           int64
	CreatedAt      time.Time
	LastModifiedAt time.Time
}

func NewFileMD(filename string) *FileMD {
	if filename == "" {
		return nil
	}
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return nil
	}
	stat := fileinfo.Sys().(*syscall.Stat_t)
	createdAt := time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)

	return &FileMD{
		Size:           fileinfo.Size(),
		CreatedAt:      createdAt,
		LastModifiedAt: fileinfo.ModTime(),
	}
}
