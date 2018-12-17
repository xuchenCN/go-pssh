package cmd

import (
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"sync"
)

func NewPsshCommand() *cobra.Command {

	config := config{}

	command := &cobra.Command{
		Use:   "pssh",
		Short: "Parallel ssh tools written in golang",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.validate(); err != nil {
				return err
			}
			return runCmd(cmd, args, config)
		},
	}

	config.addFlags(command.Flags())

	return command
}

func runCmd(cmd *cobra.Command, args []string, config config) error {

	sshWorkers, err := config.buildWorkers()

	if err != nil {
		return err
	}

	if len(sshWorkers) <= 0 {
		return fmt.Errorf("no hosts to connects")
	}

	if config.Async {
		wg := sync.WaitGroup{}
		wg.Add(len(sshWorkers))

		for _, worker := range sshWorkers {
			go func(worker sshWorker) {
				if err := worker.executeAndOutput(config.Cmd, os.Stdout, os.Stderr); err != nil {
					colorOut := color.New(color.FgRed)
					colorOut.Fprintf(os.Stderr, "[%s ERROR]\n", worker.addr)
					fmt.Fprintf(os.Stderr, "%s\n", err)
				}
				wg.Done()
			}(worker)
		}

		log.Infof("Waiting for %v hosts commands", len(sshWorkers))
		wg.Wait()

	} else {
		for _, worker := range sshWorkers {

			if err := worker.executeAndOutput(config.Cmd, os.Stdout, os.Stderr); err != nil {
				fmt.Fprintf(os.Stderr, "[%s Err]\n %v\n", worker.addr, err)
			}
		}
	}

	return nil
}

func runCmdAsync(cmd *cobra.Command, args []string, config config) error {
	sshWorkers, err := config.buildWorkers()

	if err != nil {
		return err
	}

	if len(sshWorkers) <= 0 {
		return fmt.Errorf("no hosts to connects")
	}

	return nil
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
