package cmd

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"io"
	"net"
	"os"
)

type config struct {
	hostFile string
	hosts []string
	port int
	user string
	password string
	cmd string
}

func (c *config) addFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&c.hostFile,"file","f","","file path of hosts")
	fs.StringArrayVarP(&c.hosts,"list","l",[]string{},"hosts:ip1,ip2")
	fs.IntVarP(&c.port,"port","p",22,"port of ssh connect to")
	fs.StringVarP(&c.user,"user","u","root","user")
	fs.StringVarP(&c.password,"password","P","","password")
	fs.StringVarP(&c.cmd,"cmd","c","","command")
}

func (c *config) validate() error {

	if len(c.hostFile) <= 0 && len(c.hosts) <= 0 {
		return fmt.Errorf("provide file of host or hosts list")
	}

	//if len(c.cmd) <= 0 {
	//	return fmt.Errorf("where is your command")
	//}

	return nil
}

func (c *config) listHosts() []string {
	result := make(map[string]interface{})
	if len(c.hostFile) > 0 {
		if file, err := os.Open(c.hostFile); err != nil {
			defer file.Close()
			fr := bufio.NewReader(file)
			for {
				b, _ , err := fr.ReadLine()
				if err == io.EOF {
					break;
				}
				line := string(b)
				if ip := net.ParseIP(line); ip == nil {
					log.Error("%s is not valid ip addr ignore it")
					continue
				}

				result[line] = nil
			}
		}
	}

	if len(c.hosts) > 0{
		for _,host := range c.hosts {
			if ip := net.ParseIP(host); ip == nil {
				log.Error("%s is not valid ip addr ignore it")
				continue
			}

			result[host] = nil
		}
	}

	keys := make([]string, len(result))
	for k,_ := range result {

		keys = append(keys,k)
	}

	return keys
}