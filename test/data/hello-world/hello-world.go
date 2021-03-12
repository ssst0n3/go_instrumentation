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
	pkg.NotMainPkg("main")
	log.Println("test")
	fmt.Printf("hello %s\n", name)
}

func main() {
	sayHello("hello")
	A{}.Do()
}
