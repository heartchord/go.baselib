package main

import (
	"fmt"

	"github.com/heartchord/goblazer"
)

func main() {
	if goblazer.IsFileExisted("D:\\WorkSpace") {
		fmt.Println("existed")
	} else {
		fmt.Println("not existed")
	}
}
