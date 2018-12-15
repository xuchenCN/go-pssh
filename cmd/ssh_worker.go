package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

type sshWorker struct {
	addr string
	sshClient *ssh.Client
}

func (sw *sshWorker) execute(cmd string) error {
	sess , err := sw.sshClient.NewSession()
	if err != nil {
		return nil;
	}
	defer sess.Close()

	result ,err := sess.CombinedOutput(cmd)

	suffix := "OK"
	if err != nil {
		//return err;
		suffix = "Err"
	}

	fmt.Fprintf(os.Stdout,"[%s %s] %s",sw.addr,suffix,string(result));

	return nil
}

func (sw *sshWorker) close() {
	if sw.sshClient != nil {
		sw.sshClient.Close()
	}
}