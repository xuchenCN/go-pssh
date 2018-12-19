package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
)

type sshWorker struct {
	HostSpec
	sshClient *ssh.Client
}

func (sw *sshWorker) open() error {
	conn, err := net.Dial("tcp", sw.Addr)
	if err != nil {
		return err
	}

	auth := []ssh.AuthMethod{ssh.Password(sw.Password)}

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
