package main

import (
	"strings"
)

func main1() {
	a := "asdOxyz"
	b := a[:strings.Index(a, "O")]
	println(b)
}
