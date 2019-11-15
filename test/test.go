package main

import (
	"fmt"
	"strings"
)

type A struct {
	Str string
	Arr []string
}

func main() {
	stru := A{Str: `hello`}
	a := `   123123   12312312 123123 asdfads    fgasd`
	b := strings.Fields(a)
	delSame(&b)
	fmt.Print(len(stru.Arr))
}
func delSame(arr *[]string) {
	ar := *arr
a:
	for {
		for i, n := range ar {
			if i != len(ar)-1 {
				for _, m := range ar[i+1:] {
					if n != m {
						continue
					} else {
						ar = append(ar[i+1:])
						continue a
					}
				}
			}
			break a
		}
	}
	arr = &ar
}
