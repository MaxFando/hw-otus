package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	greeting := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(greeting)
}
