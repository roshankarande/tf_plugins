package vcfile

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// src is a dir
func CreateTar(src, dst string) error {
	dir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer dir.Close()

	// get list of files
	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	// create tar file
	tarfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			return err
		}

		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = filepath.Base(file.Name())
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(tarfileWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
}

