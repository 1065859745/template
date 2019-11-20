package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

type Call struct {
	messages string
	all, to  []string
}

var call Call

func main() {
	flag.Parse()
	scan := bufio.NewScanner(os.Stdin)
	reg, _ := regexp.Compile(`^\s*$`)

	go serve()

	for {
		scan.Scan()
		t := scan.Text()
		if !reg.MatchString(t) {
			if strings.HasPrefix(t, `$Call`) {
				call.to = strings.Fields(t)[1:]
				delSame(&call.to, `127.0.0.1`)
				continue
			}
			if strings.HasPrefix(t, `$Clear all`) {
				call.all = []string{}
				continue
			}
			if strings.HasPrefix(t, `$List`) {
				fmt.Println(`-------------------------`)
				for i, v := range call.all {
					fmt.Println(i+1, v)
				}
				fmt.Println(`-------------------------`)
				continue
			}
			if strings.HasPrefix(t, `$All`) {
				call.to = []string{}
				continue
			}

			call.messages = t
			if len(call.to) != 0 {
				go call.Send(call.to)
				continue
			}
			go call.Send(call.all)
		}
	}
}

func serve() {
	ipAddr, e := net.ResolveIPAddr(`ip4`, `127.0.0.1`)
	if e != nil {
		log.Fatal(e)
	}
	conn, _ := net.ListenIP("ip4:ip", ipAddr)

	var remoteAddr string
	go func(conn *net.IPConn) {
		for {
			var b [512]byte
			n, addr, e := conn.ReadFromIP(b[:])
			if e != nil {
				log.Printf("| %s", e.Error())
				continue
			}
			remoteAddr = addr.String()
			if remoteAddr != `127.0.0.1` {
				log.Printf("|From| %s:%s", remoteAddr, b[:n])
				update(&call.all, remoteAddr)
			}
		}
	}(conn)
}

func (c Call) Send(a []string) {
	for _, v := range a {
		conn, e := net.Dial("tcp4:1", v)
		if e != nil {
			log.Printf("| %s", e.Error())
			continue
		}

		log.Printf("|To| %s", c.messages)
		conn.Write([]byte(fmt.Sprintf(" %s", c.messages)))
		conn.Close()

		if v != `127.0.0.1` {
			update(&c.all, v)
		}
	}
}

// func delNearby(arr *[]string) {
// 	ar := *arr
// a:
// 	for {
// 		for i, v := range ar {
// 			if i != len(ar)-1 {
// 				if v == ar[i+1] {
// 					ar = append(ar[:i], ar[i+1:]...)
// 					continue a
// 				}
// 			}
// 		}
// 		break
// 	}
// 	arr = &ar
// }

func delSame(ar *[]string, s string) {
	arr := *ar
a:
	for {
		for i, n := range arr {
			if i != len(arr)-1 {
				for j, m := range arr[i+1:] {
					if n != m {
						continue
					} else {
						arr = append(arr[:j+i+1], arr[j+i+2:]...)
						continue a
					}
				}
			}
		}
		break
	}
	for i, v := range arr {
		if v != s {
			continue
		}
		arr = append(arr[:i], arr[i+1:]...)
		break
	}
	*ar = arr
}

func update(arr *[]string, s string) {
	a := *arr
	for i, v := range a {
		if i != len(a)-1 {
			if v != s {
				continue
			}
			break
		}
		if v != s {
			a = append(a, s)
		}
	}
	*arr = a
}
