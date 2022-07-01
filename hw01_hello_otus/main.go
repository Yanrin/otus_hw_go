package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const str = "Hello, OTUS!"

func main() {
	reversed := stringutil.Reverse(str)
	fmt.Println(reversed)
}
