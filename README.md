# go instrumentation tool

## usage

### trace demo

```shell
╭─st0n3@yoga in ~/pentest_project/go_instrumentation on main ✘ (origin/main)
╰$ cd cmd/tracer 
╭─st0n3@yoga in ~/pentest_project/go_instrumentation/cmd/tracer on main ✘ (origin/main)
╰$ go build

╭─st0n3@yoga in ~/pentest_project/go_instrumentation/test/data/hello-world on main ✘ (origin/main)
╰$ ./test.sh 
+ go clean -cache
+ go build -work -p 1 -a -toolexec /home/st0n3/pentest_project/go_instrumentation/cmd/tracer/tracer
WORK=/tmp/go-build2421763501
╭─st0n3@yoga in ~/pentest_project/go_instrumentation/test/data/hello-world on main ✘ (origin/main)
╰$ ./hello-world 
2021/03/11 11:21:47 [TRACE] main()
2021/03/11 11:21:47 [TRACE] sayHello( name=%!v(MISSING) )
2021/03/11 11:21:47 [TRACE] NotMainPkg( from=%!v(MISSING) )
from: main
2021/03/11 11:21:47 test
hello hello
```

before instrument

```go
package main

import (
	"fmt"
	"hello-world/pkg"
	"log"
)

func debug(arg1 []string, arg2 []string) {
	fmt.Printf("name list is %v, %v\n", arg1, arg2)
}

func sayHello(name string) {
	pkg.NotMain("main")
	log.Println("test")
	fmt.Printf("hello %s\n", name)
}

func main() {
	sayHello("hello")
}
```

after instrument

```go
package main

import (
	"fmt"
	"hello-world/pkg"
	"log"
)

func debug(arg1 []string, arg2 []string) {
	{
		fmt.Printf("[TRACE] debug( arg1=%v  arg2=%v )\n", arg1, arg2)
	}
	fmt.Printf("name list is %v, %v\n", arg1, arg2)
}

func sayHello(name string) {
	{
		fmt.Printf("[TRACE] sayHello( name=%v )\n", name)
	}
	pkg.NotMain("main")
	log.Println("test")
	fmt.Printf("hello %s\n", name)
}

func main() {
	{
		fmt.Printf("[TRACE] main()\n")
	}
	sayHello("hello")
}
```

## related project

* https://github.com/ssst0n3/docker_instrumentation 