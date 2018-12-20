go-pssh is a parallel ssh tool written in golang, can define common or particular host configuration, use to deploy distributed system or execute some distributed programs like (tensorflow)[tensorflow.org]

### Example command
```
go-pssh -l <ip>,<ip> -u <user> -p <port> -P <password> -c "<command>"
```

### Example use host list file
host.txt
```
<ip>
<ip>
<ip>
...
```

```
go-pssh -f host.txt -u <user> -p <port> -P <password> -c "<command>"
```

### Example use yaml
Use same ```port user passowrd command```

yaml
```
hosts:
  - xxx.xxx.xxx.xxx
  - xxx.xxx.xxx.xxx
  - xxx.xxx.xxx.xxx
port: <port>
user: <user>
password: <password>
cmd: <command>
```

```
go-pssh -y xxxx.yaml
```

### Example host special configuration
Custome host config and merge by common configuration

yaml
```
hosts:
  - 10.110.110.12
  - 10.110.110.78
spec:
  - addr: 10.110.110.92
    user: <special user name>
  - addr: 192.168.1.134
    password: <special password>
  - addr: 10.110.110.123:<special port>
    cmd: <special command>
port: <common port>
user: <common user>
password: <common password>
cmd: <common command>
```

example.yaml
```
hosts:
  - 10.110.110.12
  - 10.110.110.78
spec:
  - addr: 10.110.110.78
    user: bar
  - addr: 192.168.1.134
    password: foo
  - addr: 10.110.110.123:10022
    cmd: "echo hello"
port: 22
user: root
password: 123456
cmd: "uname -a"
```

after merge
```
- addr: 10.110.110.12:22
  user: root
  password: 123456
  cmd: "uname -a"
- addr: 10.110.110.78:22
  user: bar
  password: 123456
  cmd: "uname -a"
- addr: 192.168.1.134:22
  user: root
  password: foo
  cmd: "uname -a"
- addr: 10.110.110.123:10022
  user: root
  password: 123456
  cmd: "echo hello"
```
