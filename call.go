package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var port = flag.String("p", "", "Local serve port, default random.")
var messages, agent string
var to []string

func main() {
	flag.Parse()
	scan := bufio.NewScanner(os.Stdin)
	reg, _ := regexp.Compile(`^\s*$`)

	if *port == `` {
		for {
			rand.Seed(time.Now().Unix())
			*port = strconv.Itoa(rand.Intn(59000) + 1000)
			tcpAddress, _ := net.ResolveTCPAddr("tcp4", ":"+*port)
			l, e := net.ListenTCP("tcp", tcpAddress)
			if e != nil {
				continue
			} else {
				l.Close()
			}
			fmt.Printf("Server start at 127.0.0.1:%s\n", *port)
			break
		}
	}

	go serve()

	for {
		scan.Scan()
		t := scan.Text()
		if !reg.MatchString(t) {
			if strings.HasPrefix(t, `$Call`) {
				to = strings.Fields(t)[1:]
				delSame(&to)
			} else if strings.HasPrefix(t, `$Agent`) {

			} else {
				messages = t
				_, e := http.Get(fmt.Sprintf("http://127.0.0.1:%s/call", *port))
				if e != nil {
					log.Printf("| %s", e.Error())
				} else {
					log.Printf("| %s", messages)
				}
			}
		}
	}
}

func serve() {
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		remoteIp := strings.Split(r.RemoteAddr, `:`)[0]
		if remoteIp == `localhost` || remoteIp == `127.0.0.1` {
			for _, v := range to {
				if _, e := http.Get(fmt.Sprintf("http://%s/conversation?messages=%s", v, fmt.Sprintf("%s%%20:%%20%s", *port, messages))); e != nil {
					log.Printf("| Faild to call %s", v)
				}
			}
		}
	})
	http.HandleFunc("/conversation", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if m, ok := query["messages"]; ok {
			log.Printf("| %s:%s", strings.Split(r.RemoteAddr, `:`)[0], m[0])
		}
		fmt.Fprint(w, `Ok`)
	})
	http.ListenAndServe(`:`+*port, nil)
}

func delNearby(arr *[]string) {
	ar := *arr
a:
	for {
		for i, v := range ar {
			if i != len(ar)-1 {
				if v == ar[i+1] {
					ar = append(ar[:i], ar[i+1:]...)
					continue a
				}
			}
		}
		break
	}
	arr = &ar
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
