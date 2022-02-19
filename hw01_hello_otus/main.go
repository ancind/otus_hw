package main

import "fmt"
import "golang.org/x/example/stringutil"

func main() {
	str := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(str)
}
