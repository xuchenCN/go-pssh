package cmd

import (
	"io"
	"os"
)

func CopyFile(srcFile io.Reader, dstFile io.Writer) (int64, error) {
	return io.Copy(dstFile, srcFile)
}

func IsDir(file *os.File) (bool, error) {
	fi, err := file.Stat()
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}
