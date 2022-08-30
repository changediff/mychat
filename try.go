package main

import (
	"fmt"
)

func main() {
	msg := "1 2 3 4 5 6   7"
	fmt.Println(Colorize(msg, FgRed, BgBlue))
	// fmt.Println(msg)
}
