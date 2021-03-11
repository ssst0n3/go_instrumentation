package pkg

import "fmt"

func NotMainPkg(from string) {
	fmt.Printf("from: %s\n", from)
}
