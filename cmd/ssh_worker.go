package cmd

import (
	"crypto/x509"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
)

type sshWorker struct {
	HostSpec
	KeyPath      string
	KeyEncrypted bool
	sshClient    *ssh.Client
}

func (sw *sshWorker) open() error {
	conn, err := net.Dial("tcp", sw.Addr)
	if err != nil {
		return err
	}

	auth := sw.buildAuth()

	sshConf := &ssh.ClientConfig{User: sw.User, Auth: auth, HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}}

	sshConn, newChan, reqChan, err := ssh.NewClientConn(conn, sw.Addr, sshConf)
	if err != nil {
		return err
	}

	sw.sshClient = ssh.NewClient(sshConn, newChan, reqChan)

	return nil
}

func (sw *sshWorker) buildAuth() []ssh.AuthMethod {
	result := []ssh.AuthMethod{}
	if len(sw.Password) > 0 {
		result = append(result, ssh.Password(sw.Password))
	}

	if len(sw.KeyPath) > 0 {
		key, err := ioutil.ReadFile(sw.KeyPath)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}

		var signer ssh.Signer
		if sw.KeyEncrypted {
			encryptedKey, err := x509.ParsePKCS1PrivateKey(key)
			if err != nil {
				log.Fatal(err)
			}
			signer, err = ssh.NewSignerFromKey(encryptedKey)
			if err != nil {
				log.Fatalf("unable to parse private key: %v", err)
			}
		} else {
			// Create the Signer for this private key.
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				log.Fatalf("unable to parse private key: %v", err)
			}
		}

		result = append(result, ssh.PublicKeys(signer))
	}

	return result
}

func (sw *sshWorker) remoteCopy(src, distDir string) error {

	if sw.sshClient == nil {
		if err := sw.open(); err != nil {
			return err
		}
	}

	if !path.IsAbs(distDir) {
		return fmt.Errorf(distDir + " is not absolute path")
	}

	if !path.IsAbs(src) {
		return fmt.Errorf(src + " is not absolute path")
	}

	sftpClient, err := sftp.NewClient(sw.sshClient)

	if err != nil {
		log.Error("Errot to create sftp")
		return err
	}

	defer sftpClient.Close()

	srcFile, err := os.Open(src)

	if err != nil {
		return err
	}
	defer srcFile.Close()

	if isDir, err := IsDir(srcFile); err != nil {
		return err
	} else if isDir { // Directory
		if err := copyDir(sftpClient, src, path.Join(distDir, path.Base(src))); err != nil {
			return err
		}
	} else { // File
		if _, err := copyFile(sftpClient, src, distDir); err != nil {
			return err
		}
	}

	log.Infof("%s copy done! \n", sw.Addr)

	return nil
}

func copyFile(sftpClient *sftp.Client, src, distDir string) (int64, error) {

	fileName := filepath.Base(src)

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	srcStat, _ := os.Stat(src)

	dstFile, err := sftpClient.Create(path.Join(distDir, fileName))
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	sftpClient.Chmod(dstFile.Name(), srcStat.Mode())

	return CopyFile(srcFile, dstFile)

}

func copyDir(sftpClient *sftp.Client, src, distDir string) error {

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return nil
	}

	srcStat, _ := os.Stat(src)

	err = sftpClient.Mkdir(distDir)
	if err != nil && !os.IsExist(err) {
		return err
	}

	sftpClient.Chmod(distDir, srcStat.Mode())

	for _, file := range files {
		if file.IsDir() {
			if err = copyDir(sftpClient, path.Join(src, file.Name()), path.Join(distDir, file.Name())); err != nil {
				return err
			}
		} else {
			if _, err = copyFile(sftpClient, path.Join(src, file.Name()), distDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func (sw *sshWorker) executeAndOutput(stdout io.Writer, stderr io.Writer) error {

	if sw.sshClient == nil {
		if err := sw.open(); err != nil {
			return err
		}
	}

	defer sw.close()

	sess, err := sw.sshClient.NewSession()
	if err != nil {
		return nil
	}

	defer sess.Close()

	result, err := sess.CombinedOutput(sw.Cmd)

	if err != nil {
		colorOut := color.New(color.FgRed)
		colorOut.Fprintf(stderr, "[%s %s]\n ", sw.Addr, "ERROR")
		fmt.Fprintf(stderr, "%s %s\n", string(result), err)
	} else {
		colorOut := color.New(color.FgGreen)
		colorOut.Fprintf(stdout, "[%s %s]\n", sw.Addr, "OK")
		fmt.Fprintf(stdout, "%s\n", string(result))
	}

	return nil
}

func (sw *sshWorker) close() {
	if sw.sshClient != nil {
		sw.sshClient.Close()
	}
}
