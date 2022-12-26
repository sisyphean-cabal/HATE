package main

import "fmt"

func main() {
	fmt.Println("Hello, Arch!")

	type Options struct {
		Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	}
}
