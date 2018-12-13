package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

func NewPsshCommand() *cobra.Command {

	config := config{}

	command := &cobra.Command{
		Use: "pssh",
		Short: "Parallel ssh tools written in golang",
		RunE:func(cmd *cobra.Command, args []string) error {
			if err := config.validate(); err != nil {
				return err;
			}
			return runCmd(cmd,args,config)
		},
	};

	config.addFlags(command.Flags())

	return command;
}

func runCmd(cmd *cobra.Command, args []string, config config) error {

	log.Info("Execute", config)

	hosts := config.listHosts()

	if len(hosts) <= 0 {
		return fmt.Errorf("no hosts to connects")
	}

	sshWorkers := make(map[string]sshWorker,len(hosts))

	for host := range hosts {
		addr := fmt.Sprintf("%s:%v",host,config.port)
		conn , err := net.Dial("tcp",addr)
		if err != nil {
			return err
		}

		auth := []ssh.AuthMethod{ssh.Password(config.password)}

		sshConf := &ssh.ClientConfig{User:config.user,Auth:auth}
		sshConn, newChan, reqChan, err := ssh.NewClientConn(conn,addr,sshConf)
		if err != nil {
			return err
		}

		sshClient := ssh.NewClient(sshConn,newChan,reqChan)
		sshWorker := sshWorker{sshClient:sshClient}
		sshWorkers[addr] = sshWorker
	}



	return nil;
}

func ReadKey(sshkey string, privateKey *[]ssh.AuthMethod) bool {

	buf, err := ioutil.ReadFile(sshkey)
	if err != nil {
		return false
	}
	signer, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return false
	}
	*privateKey = append(*privateKey, ssh.PublicKeys(signer))
	return true
}

