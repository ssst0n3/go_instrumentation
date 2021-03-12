package main

import "fmt"

type A struct {
}

func (a A) Do() {
	fmt.Println("a.Do()")
}
