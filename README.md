Parallel ssh tool written in golang

### Example
```
go-pssh -l xxx.xxx.xxx.xxx,xxx.xxx.xxx.xxx -u user -p 22 -P password -c "command"
```

### Example use yaml

yaml
```
hosts:
  - xxx.xxx.xxx.xxx
  - xxx.xxx.xxx.xxx
  - xxx.xxx.xxx.xxx
port: 22
user: user
password: password
cmd: uname -a
```

```
go-pssh -y xxxx.yaml
```
