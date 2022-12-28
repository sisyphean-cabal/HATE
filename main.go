package main

import (
	"flag"
	"fmt"
)

func main() {
	var hatectl bool

	flag.BoolVar(&hatectl, "hctl", false, "WHY ARE YOU DOING THIS")
	flag.Parse()
	fmt.Println(hatectl)
}
