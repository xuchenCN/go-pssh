package cmd

import (
	"io"
	"os"
)

func CopyFile(srcFile io.Reader , dstFile io.Writer , bufLen int32) error {

	if bufLen <= 0 {
		bufLen = 2048
	}

	buf := make([]byte, bufLen)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	return nil
}

func IsDir(file *os.File) (bool, error) {
	fi ,err := file.Stat()
	if err != nil {
		return false,err
	}

	return fi.IsDir(), nil
}
