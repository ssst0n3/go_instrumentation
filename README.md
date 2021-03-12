# go instrumentation tool

## usage

### trace runc
```shell
╭─st0n3@yoga in ~/pentest_target 
╰$ git clone https://github.com/opencontainers/runc
正克隆到 'runc'...
remote: Enumerating objects: 3, done.
remote: Counting objects: 100% (3/3), done.
remote: Compressing objects: 100% (3/3), done.
remote: Total 27503 (delta 0), reused 3 (delta 0), pack-reused 27500
接收对象中: 100% (27503/27503), 11.60 MiB | 109.00 KiB/s, 完成.
处理 delta 中: 100% (17980/17980), 完成.
╭─st0n3@yoga in ~/pentest_target 
╰$ cd runc 
╭─st0n3@yoga in ~/pentest_target/runc on master ✔ (origin/master)
╰$ EXTRA_FLAGS="-a -toolexec /home/st0n3/pentest_project/go_instrumentation/cmd/tracer/tracer" make 
go build -trimpath "-buildmode=pie" -a -toolexec /home/st0n3/pentest_project/go_instrumentation/cmd/tracer/tracer -tags "seccomp" -ldflags "-X main.gitCommit="249bca0a1316129dcd5bd38b5d75572274181cb5" -X main.version=1.0.0-rc93+dev " -o runc .
╭─st0n3@yoga in ~/pentest_target/runc on master ✔ (origin/master)
╰$ ./runc --version
runc version 1.0.0-rc93+dev
commit: 249bca0a1316129dcd5bd38b5d75572274181cb5
spec: 1.0.2-dev
go: go1.16
libseccomp: 2.5.1
╭─st0n3@yoga in ~/pentest_target/runc on master ✔ (origin/master)
╰$ cat /tmp/instrumentation 
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/godbus/dbus/v5 init()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/godbus/dbus/v5 detectEndianness()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/godbus/dbus/v5 init()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/godbus/dbus/v5 init()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/godbus/dbus/v5 init()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/pkg/errors Wrap( err=too many levels of symbolic links  message=secure join )"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/pkg/errors callers()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/opencontainers/runc/libcontainer/utils init()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/opencontainers/runc/libcontainer/cgroups GetHugePageSize()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/opencontainers/runc/libcontainer/cgroups getHugePageSizeFromFilenames( fileNames=[hugepages-2048kB hugepages-1048576kB] )"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/pkg/errors New( message=cgroup: subsystem does not exist )"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/pkg/errors callers()"
time="2021-03-12T11:32:09+08:00" level=info msg="[TRACE] github.com/cilium/ebpf/internal init()"
...
```

### trace demo

```shell
╭─st0n3@vmhome in ~/pentest_project/go_instrumentation on main ✔ (origin/main)
╰$ cd cmd/tracer 
╭─st0n3@vmhome in ~/pentest_project/go_instrumentation/cmd/tracer on main ✔ (origin/main)
╰$ ./build.sh 
╭─st0n3@vmhome in ~/pentest_project/go_instrumentation/cmd/tracer on main ✔ (origin/main)
╰$ cd ../../test/data/hello-world/   
╭─st0n3@vmhome in ~/pentest_project/go_instrumentation/test/data/hello-world on main ✔ (origin/main)
╰$ ./test.sh 
WORK=/tmp/go-build823317741
╭─st0n3@vmhome in ~/pentest_project/go_instrumentation/test/data/hello-world on main ✔ (origin/main)
╰$ ./hello-world                                      
from: main
pkg2
2021/03/12 01:44:51 test
hello hello
╭─st0n3@vmhome in ~/pentest_project/go_instrumentation/test/data/hello-world on main ✔ (origin/main)
╰$ cat /tmp/instrumentation 
time="2021-03-12T01:44:51+08:00" level=info msg="[TRACE] main main()\n"
time="2021-03-12T01:44:51+08:00" level=info msg="[TRACE] main sayHello( name=hello )\n"
time="2021-03-12T01:44:51+08:00" level=info msg="[TRACE] hello-world/pkg NotMainPkg( from=main )\n"
time="2021-03-12T01:44:51+08:00" level=info msg="[TRACE] hello-world/pkg2 Log()\n"
```

before instrument

```go
package main

import (
	"fmt"
	"hello-world/pkg"
	"hello-world/pkg2"
	"log"
)

func debug(arg1 []string, arg2 []string) {
	fmt.Printf("name list is %v, %v\n", arg1, arg2)
}

func sayHello(name string) {
	pkg.NotMainPkg("main")
	pkg2.Log()
	log.Println("test")
	fmt.Printf("hello %s\n", name)
}

func main() {
	sayHello("hello")
}
```

after instrument

```go
//line /home/st0n3/pentest_project/go_instrumentation/test/data/hello-world/hello-world.go:1
package main

import (
	"fmt"
	instrument_os "os"
	"hello-world/pkg"
	"hello-world/pkg2"
	"log"
	instrument_log "github.com/sirupsen/logrus"
)

func debug(arg1 []string, arg2 []string) {
	{
		logger := instrument_log.New()
		file, _ := instrument_os.OpenFile("/tmp/instrumentation", instrument_os.O_CREATE|instrument_os.O_WRONLY|instrument_os.O_APPEND, 0644)
		logger.SetOutput(file)
		logger.Infof("[TRACE] main debug( arg1=%v  arg2=%v )\n", arg1, arg2)
	}
	fmt.Printf("name list is %v, %v\n", arg1, arg2)
}

func sayHello(name string) {
	{
		logger := instrument_log.New()
		file, _ := instrument_os.OpenFile("/tmp/instrumentation", instrument_os.O_CREATE|instrument_os.O_WRONLY|instrument_os.O_APPEND, 0644)
		logger.SetOutput(file)
		logger.Infof("[TRACE] main sayHello( name=%v )\n", name)
	}
	pkg.NotMainPkg("main")
	pkg2.Log()
	log.Println("test")
	fmt.Printf("hello %s\n", name)
}

func main() {
	{
		logger := instrument_log.New()
		file, _ := instrument_os.OpenFile("/tmp/instrumentation", instrument_os.O_CREATE|instrument_os.O_WRONLY|instrument_os.O_APPEND, 0644)
		logger.SetOutput(file)
		logger.Infof("[TRACE] main main()\n")
	}
	sayHello("hello")
}
```

## related project

* https://github.com/ssst0n3/docker_instrumentation 

