package rz

import "fmt"

type quiet bool

func (q quiet) println(args ...interface{}) {
	if !q {
		fmt.Println(args...)
	}
}

func passed(product int, filter string) bool {
	switch filter {
	case "odd":
		return product&1 != 0
	case "even":
		return product&1 == 0
	case "any":
		return true
	default:
		return false
	}
}

