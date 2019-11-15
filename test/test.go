package main

import (
	"fmt"
	"net/http"
	"strings"
)

type A struct {
	name string
}

func main() {
	str := A{name: `zhang`}
	str.b()
	fmt.Println(str.name)
}

func (a A) b() {
	a.name = `wang`
	fmt.Println(a.name)
}
func serve() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		if r.Host == `127.0.0.1` || r.Host == `localhost` {
			query := r.URL.Query()
			if to, ok := query["to"]; ok {
				to = strings.Split(to[0], `,`)
				for _, v := range to {
					fmt.Println(v)

				}
			}
		}
		fmt.Fprint(w, `hello`)
	})

	err := http.ListenAndServe(`127.0.0.1:80`, nil)
	if err != nil {
		fmt.Println(`err`)
	}
	return
}
