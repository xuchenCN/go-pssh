package cmd

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/xuchenCN/go-pssh/yaml"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Kubelet component name
	scpCommand = "scp"
)

type config struct {
	hostFile   string
	hostList   string
	configFile string

	Hosts        []string   `json:"hosts"`
	Port         int        `json:"port"`
	User         string     `json:"user"`
	Password     string     `json:"password"`
	KeyPath      string     `json:"key_path"`
	KeyEncrypted bool       `json:"key_encrypted"`
	Cmd          string     `json:"cmd"`
	Async        bool       `json:"async"`
	HostSpec     []HostSpec `json:"spec"`

	scpLocal  string
	scpRemote string
}

type HostSpec struct {
	Addr     string `json:"addr"`
	Cmd      string `json:"cmd"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (c *config) addFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&c.configFile, "config", "y", "", "config file format in yaml or json")
	fs.StringVarP(&c.hostFile, "file", "f", "", "file path of hosts")
	fs.StringVarP(&c.hostList, "list", "l", "", "hosts:ip1,ip2")
	fs.BoolVarP(&c.Async, "async", "a", false, "execute concurrently")
	fs.BoolVarP(&c.KeyEncrypted, "key-encrypted", "e", false, "encrypted private key")
	fs.IntVarP(&c.Port, "port", "p", 22, "port of ssh connect to")
	fs.StringVarP(&c.User, "user", "u", "root", "user")
	fs.StringVarP(&c.Password, "password", "P", "", "password")
	fs.StringVarP(&c.KeyPath, "key", "k", "", "private key")
	fs.StringVarP(&c.Cmd, "cmd", "c", "", "command")

	//Scp command
	fs.StringVarP(&c.scpLocal, "src", "s", "", "local file or directory when scp")
	fs.StringVarP(&c.scpRemote, "dist", "d", "", "remote directory when scp")
}

func (c *config) validate(subCmd string) error {

	if len(c.hostFile) <= 0 && len(c.hostList) <= 0 && len(c.configFile) <= 0 {
		return fmt.Errorf("need file of host or hosts list or config file")
	}

	if len(c.configFile) > 0 {
		if abs, err := filepath.Abs(c.configFile); err != nil {
			return nil
		} else {
			c.configFile = abs
			log.Infof("convert path to %s", c.configFile)
		}

		c.loadConfigFile()
	}

	if len(c.Cmd) <= 0 && subCmd != scpCommand {
		return fmt.Errorf("where is your command")
	}

	if len(c.hostFile) > 0 {
		if abs, err := filepath.Abs(c.hostFile); err != nil {
			return nil
		} else {
			c.hostFile = abs
			log.Infof("convert path to %s", c.hostFile)
		}
	}

	if len(subCmd) > 0 {

		switch subCmd {
		case scpCommand:

			if len(c.scpLocal) <= 0 {
				return fmt.Errorf("Using -s to specify local file or directory when scp")
			}

			if len(c.scpRemote) <= 0 {
				return fmt.Errorf("Using -d to specify remote directory when scp")
			}

			if len(c.scpLocal) > 0 {
				if abs, err := filepath.Abs(c.scpLocal); err != nil {
					return nil
				} else {
					c.scpLocal = abs
					log.Infof("convert path to %s", c.scpLocal)
				}
			}
			break
		}

	}

	return nil
}

func (c *config) loadConfigFile() error {
	if len(c.configFile) <= 0 {
		return fmt.Errorf("use -cfg to locate config-file")
	}

	if !filepath.IsAbs(c.configFile) {
		return fmt.Errorf("expect the abosulte path of config-file")
	}

	if cfgFile, err := os.Open(c.configFile); err == nil {
		yamlToJsonDecoder := yaml.NewYAMLToJSONDecoder(cfgFile)
		return yamlToJsonDecoder.Decode(&c)
	} else {
		return err
	}
}

func (c *config) listHosts() ([]string, error) {
	result := make(map[string]interface{})
	if len(c.hostFile) > 0 {
		if file, err := os.Open(c.hostFile); err == nil {
			defer file.Close()
			fr := bufio.NewReader(file)
			for {
				b, _, err := fr.ReadLine()
				if err == io.EOF {
					break
				}
				line := strings.TrimSpace(string(b))
				if ip := net.ParseIP(line); ip == nil {
					log.Error("%s is not valid ip addr ignore it")
					continue
				}

				result[line] = nil
			}
		} else {
			return nil, err
		}
	}

	if len(c.hostList) > 0 {
		list := strings.Split(c.hostList, ",")
		for _, host := range list {
			host = strings.TrimSpace(host)
			if ip := net.ParseIP(host); ip == nil {
				log.Error("%s is not valid ip addr ignore it")
				continue
			}

			result[host] = nil
		}
	}

	if len(c.Hosts) > 0 {
		for _, host := range c.Hosts {
			result[host] = nil
		}
	}

	keys := make([]string, 0, len(result))
	for k, _ := range result {
		keys = append(keys, k)
	}

	return keys, nil
}

func (c *config) buildWorkers() (map[string]sshWorker, error) {

	hosts, err := c.listHosts()

	if err != nil {
		return nil, err
	}

	sshWorkers := make(map[string]sshWorker, len(hosts))

	for _, host := range hosts {
		addr := fmt.Sprintf("%s:%v", host, c.Port)

		sshWorker := sshWorker{HostSpec: HostSpec{User: c.User,
			Addr:     addr,
			Password: c.Password,
			Cmd:      c.Cmd},
			KeyPath:      c.KeyPath,
			KeyEncrypted: c.KeyEncrypted,
		}

		sshWorkers[host] = sshWorker
	}

	//Load host special configuration
	if len(c.HostSpec) > 0 {
		for _, spec := range c.HostSpec {

			addrSplited := strings.Split(spec.Addr, ":")
			host := addrSplited[0]

			worker, ok := sshWorkers[host]
			if !ok {
				worker = sshWorker{HostSpec: spec}
			}

			//User
			if len(spec.User) > 0 {
				worker.User = spec.User
			} else {
				worker.User = c.User
			}
			//Password
			if len(spec.Password) > 0 {
				worker.Password = spec.Password
			} else {
				worker.Password = c.Password
			}
			//Cmd
			if len(spec.Cmd) > 0 {
				worker.Cmd = spec.Cmd
			} else {
				worker.Cmd = c.Cmd
			}
			//Addr string has port
			if len(addrSplited) > 1 {
				worker.Addr = spec.Addr
			} else {
				worker.Addr = fmt.Sprintf("%s:%v", host, c.Port)
			}

			sshWorkers[host] = worker
		}
	}

	return sshWorkers, nil
}
