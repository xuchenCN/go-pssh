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
	fs.StringVar(&c.hostFile,"f","","file path of hosts")
	fs.StringArrayVar(&c.hosts,"h",[]string{},"hosts:ip1,ip2")
	fs.IntVar(&c.port,"port",22,"port of ssh connect to")
	fs.StringVar(&c.user,"u","root","user")
	
}
