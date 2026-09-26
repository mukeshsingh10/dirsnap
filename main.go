package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println(`usage:
  dirsnap init
  dirsnap snap
  dirsnap restore <hash> <target-dir>
  dirsnap log`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if err := os.MkdirAll(objectsDir(), 0755); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("Initialized empty dirsnap repository in .dirsnap/")

	case "snap":
		hash, err := Snap(".")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("Snapshot:", hash)

	case "restore":
		if len(os.Args) < 4 {
			usage()
			os.Exit(1)
		}
		if err := Restore(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("Restored into", os.Args[3])

	case "log":
		if err := Log(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		usage()
		os.Exit(1)
	}
}
