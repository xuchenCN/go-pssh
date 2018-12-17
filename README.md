Parallel ssh tool written in golang

### Example command
```
go-pssh -l <ip>,<ip> -u <user> -p <port> -P <password> -c "<command>"
```
host.txt
### Example use host list file
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
