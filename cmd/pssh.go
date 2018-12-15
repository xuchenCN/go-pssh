package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"os"
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

	hosts,err := config.listHosts()

	if err != nil {
		return err
	}

	log.Infof("Execute %v", hosts)

	if len(hosts) <= 0 {
		return fmt.Errorf("no hosts to connects")
	}

	sshWorkers := make(map[string]sshWorker,len(hosts))

	for _,host := range hosts {
		addr := fmt.Sprintf("%s:%v",host,config.Port)
		conn , err := net.Dial("tcp",addr)
		if err != nil {
			return err
		}

		auth := []ssh.AuthMethod{ssh.Password(config.Password)}

		sshConf := &ssh.ClientConfig{User:config.User,Auth:auth,HostKeyCallback:
			func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
		}}

		sshConn, newChan, reqChan, err := ssh.NewClientConn(conn,addr,sshConf)
		if err != nil {
			return err
		}

		sshClient := ssh.NewClient(sshConn,newChan,reqChan)
		sshWorker := sshWorker{sshClient:sshClient,addr:addr}
		sshWorkers[addr] = sshWorker
	}

	for _, worker := range sshWorkers {

		if err := worker.execute(config.Cmd); err != nil {
			fmt.Fprintf(os.Stderr,"[%s Err] %v\n",worker.addr,err)
		}
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

