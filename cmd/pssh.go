package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewPsshCommand() *cobra.Command {

	command := &cobra.Command{
		Use: "pssh",
		Short: "Parallel ssh tools written in golang",
		RunE:runCmd,
	};

	addFlags(command.Flags())

	return command;
}

func runCmd(cmd *cobra.Command, args []string) error {


	return nil;
}

