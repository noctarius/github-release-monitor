package main

import (
	"os"
	"runtime"
	"strings"
	"fmt"
)

func main() {
	for _, arg := range os.Args {
		if arg == "-a" {
			fmt.Println(runtime.GOARCH)
			return
		} else if arg == "-o" {
			fmt.Println(runtime.GOOS)
			return
		} else if arg == "-v" {
			fmt.Println(strings.Replace(runtime.Version(), "go", "", -1))
			return
		}
	}
}
