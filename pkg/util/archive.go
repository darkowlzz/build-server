package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Pack creates a tar archive and compresses a given source directory and writes
// to the given writer.
// src must be an absolute path.
// excludeDirs are the directory paths relative to the source, which are to be
// excluded from the archive.
func Pack(src string, writer io.Writer, excludeDirs []string) error {
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("failed to check stat of source: %v", err)
	}

	// Make the exclude dir paths absolute to make it easy for file name
	// comparison below.
	excludeMap := make(map[string]bool)
	for _, dir := range excludeDirs {
		excludeMap[filepath.Join(src, dir)] = true
	}

	gzw := gzip.NewWriter(writer)
	tw := tar.NewWriter(gzw)

	err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if found in exclude list
		if excludeFile(file, excludeMap) {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// Set the file name with proper path for the untar to be in proper
		// structure.
		header.Name, err = filepath.Rel(src, file)
		if err != nil {
			return err
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Directories won't have a content to write, only the header.
		if !fi.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})

	if err != nil {
		return err
	}

	tw.Close()
	gzw.Close()

	return nil
}

// excludeFile takes a file and a set of directories to exclude and checks if
// the given file is part of any of the exclude directories.
func excludeFile(file string, excludeDirs map[string]bool) bool {
	// Check if the file is one of the excluded dirs.
	if excludeDirs[file] {
		return true
	}

	// Check if the file is under an excluded dir.
	for d := range excludeDirs {
		if strings.HasPrefix(file, d) {
			return true
		}
	}

	return false
}

// Unpack reads a given reader, uncompresses the gzipped data, extracts the tar
// files and writes the files at the destination.
func Unpack(dst string, reader io.Reader) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	// Untar into files.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}
		if hdr == nil {
			continue
		}

		target := filepath.Join(dst, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		}
	}

	return nil
}
