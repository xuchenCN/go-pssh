package cmd

import "golang.org/x/crypto/ssh"

type sshWorker struct {
	sshClient *ssh.Client
}

func (sw *sshWorker) execute(cmd string) error {
	sess , err := sw.sshClient.NewSession()
	if err != nil {
		return nil;
	}

	return nil
}

func (sw *sshWorker) close() {
	if sw.sshClient != nil {
		sw.sshClient.Close()
	}
}