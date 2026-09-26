package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Snap reads every file directly inside dir (v1 is deliberately flat —
// no subdirectories yet), stores each one as an object, and then builds a
// "listing" — a plain text file mapping filename -> content hash. That
// listing is itself hashed and stored, so the ENTIRE directory's state is
// now represented by one hash. That single hash is what gets recorded as
// the snapshot.
func Snap(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var lines []string
	for _, e := range entries {
		if e.IsDir() {
			continue // v1: skip subdirectories entirely, on purpose (see README)
		}
		if e.Name() == dirsnapDir {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return "", err
		}

		fileHash, err := HashAndStore(data)
		if err != nil {
			return "", err
		}

		lines = append(lines, fmt.Sprintf("%s %s", fileHash, e.Name()))
	}

	// Sorting matters: os.ReadDir order is already alphabetical on most
	// platforms, but we sort explicitly so the listing text — and therefore
	// its hash — never depends on filesystem/OS quirks. Same files, same
	// hash, always.
	sort.Strings(lines)

	listing := strings.Join(lines, "\n")
	folderHash, err := HashAndStore([]byte(listing))
	if err != nil {
		return "", err
	}

	// The log is just an append-only text file of "when + what". We don't
	// need a separate linked-list of commit objects with parent pointers —
	// the log file's own line order already IS the history.
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\t%s\n", time.Now().Format(time.RFC3339), folderHash)
	return folderHash, err
}

// Restore takes a folder-hash (printed by Snap, or read from the log) and
// rebuilds the files it represents into targetDir. This is deliberately the
// mirror image of Snap: read the listing instead of building it, read each
// blob instead of writing it.
func Restore(folderHash, targetDir string) error {
	data, err := ReadObject(folderHash)
	if err != nil {
		return fmt.Errorf("no snapshot found for hash %s: %w", folderHash, err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	listing := strings.TrimSpace(string(data))
	if listing == "" {
		return nil // empty snapshot, nothing to restore
	}

	for _, line := range strings.Split(listing, "\n") {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		fileHash, name := parts[0], parts[1]

		content, err := ReadObject(fileHash)
		if err != nil {
			return fmt.Errorf("missing blob for %s: %w", name, err)
		}

		if err := os.WriteFile(filepath.Join(targetDir, name), content, 0644); err != nil {
			return err
		}
	}

	return nil
}

// Log prints the snapshot history as-is. No parsing needed beyond what's
// already in the file, because the file format IS the display format.
func Log() error {
	data, err := os.ReadFile(logFile)
	if err != nil {
		return fmt.Errorf("no snapshots yet")
	}
	fmt.Print(string(data))
	return nil
}
