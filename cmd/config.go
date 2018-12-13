package cmd

import "github.com/spf13/pflag"

type config struct {
	hostFile string
	hosts []string
	port int
	user string
	password string
	cmd string
}

func (c *config) addFlags(fs *pflag.FlagSet) {
	

}
