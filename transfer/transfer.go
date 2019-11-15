package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var port = flag.String("p", `80`, "Start port for transfer server.")

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if to, ok := query["to"]; ok {
			to = strings.Split(to[0], `,`)
			for _, v := range to {
				go testUrl(fmt.Sprintf("http://%s/test"))
			}
		}
	})

	fmt.Printf("Transfer server well start at http://localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func testUrl(url string, c chan bool) {
	r, e := http.Get(url)

	if e != nil {
		c <- false
		return
	}
	if r.StatusCode != 200 {
		c <- false
		return
	}
	c <- true
}
