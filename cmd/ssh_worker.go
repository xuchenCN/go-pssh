package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"io"
)

type sshWorker struct {
	addr      string
	sshClient *ssh.Client
}

func (sw *sshWorker) executeAndOutput(cmd string, stdout io.Writer, stderr io.Writer) error {
	sess, err := sw.sshClient.NewSession()
	if err != nil {
		return nil
	}
	defer sess.Close()

	result, err := sess.CombinedOutput(cmd)

	if err != nil {
		colorOut := color.New(color.FgRed)
		colorOut.Fprintf(stderr, "[%s %s]\n ", sw.addr, "ERROR")
		fmt.Fprintf(stderr, "%s %s\n", string(result), err)
	} else {
		colorOut := color.New(color.FgGreen)
		colorOut.Fprintf(stdout, "[%s %s]\n", sw.addr, "OK")
		fmt.Fprintf(stdout, "%s\n", string(result))
	}

	return nil
}

func (sw *sshWorker) close() {
	if sw.sshClient != nil {
		sw.sshClient.Close()
	}
}
