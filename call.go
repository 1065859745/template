package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var port, chStr, t, all, name = flag.String("p", "", `Local serve port, default random.`), make(chan string, 10), flag.Int("t", 0, `Start type,-1 only serve, 0 double, 1 only client.`), []string{}, `[Unkown]`

func main() {
	flag.Parse()
	scan := bufio.NewScanner(os.Stdin)
	reg, _ := regexp.Compile(`^\s*$`)
	switch *t {
	case -1:
		serve()
	case 1:
		// client()
	default:
		go serve()
		for {
			scan.Scan()
			t := scan.Text()
			if !reg.MatchString(t) {
				if strings.HasPrefix(t, `$Name`) {
					name = `[` + strings.Join(strings.Fields(t)[1:], ` `) + `]`
					continue
				}
				if strings.HasPrefix(t, `$Call`) {
					for _, v := range strings.Fields(t)[1:] {
						go client(v)
					}
					continue
				}
				if strings.HasPrefix(t, `$Clear`) {
					all = []string{}
					continue
				}
				if strings.HasPrefix(t, `$List`) {
					fmt.Println(`-------------------------`)
					for i, v := range all {
						fmt.Println(i+1, v)
					}
					fmt.Println(`-------------------------`)
					continue
				}
				if strings.HasPrefix(t, `$All`) {
					chStr <- t
					continue
				}
				if strings.HasPrefix(t, `$Answer`) {
					for _, v := range strings.Fields(t)[1:] {
						chStr <- `$Answer ` + v
					}
					continue
				}
				if strings.HasPrefix(t, `$Down`) {
					for _, v := range strings.Fields(t)[1:] {
						chStr <- `$Down ` + v
					}
					continue
				}
				t = fmt.Sprintf("%s : %s", name, t)
				chStr <- t
			}
		}
	}
}

func serve() {
	if *port == `` {
		for {
			rand.Seed(time.Now().Unix())
			*port = strconv.Itoa(rand.Intn(59000) + 1000)
			l, e := net.Listen("tcp", ":"+*port)
			if e != nil {
				continue
			} else {
				l.Close()
			}
			fmt.Printf("Serve at 127.0.0.1:%s\n", *port)
			break
		}
	}
	l, _ := net.Listen("tcp4", ":"+*port)

	for {
		conn, e := l.Accept()
		if e != nil {
			log.Printf("| %s", e.Error())
			continue
		}
		update(&all, conn.RemoteAddr().String())
		// read
		go tcpRead(conn)

		// write
		tcpWrite(conn, chStr)
	}
}

func client(addr string) {
	conn, e := net.Dial("tcp4", addr)
	if e != nil {
		log.Printf("| %s", e.Error())
		return
	}
	update(&all, addr)

	// read
	go tcpRead(conn)

	// write
	tcpWrite(conn, chStr)
}

func tcpRead(conn net.Conn) {
	var b [512]byte
	remoteAddr := conn.RemoteAddr().String()

	for {
		n, e := conn.Read(b[:])
		if e != nil {
			log.Printf("| %s", e.Error())
			log.Printf("|Miss| %s", remoteAddr)
			break
		}
		log.Printf("|From| %s %s", remoteAddr, string(b[:n]))
	}
}

func tcpWrite(conn net.Conn, chStr chan string) {
	remoteAddr := conn.RemoteAddr().String()

	for {
		m := <-chStr
		if strings.HasPrefix(m, `$Answer`) {
			if strings.HasPrefix(m, `$Answer `+remoteAddr) {
				m = <-chStr
				if !strings.HasPrefix(m, `$`) {
					_, e := conn.Write([]byte(m))
					if e != nil {
						log.Printf("| %s", e.Error())
						continue
					}
					log.Printf("| %s", m)
				}
			}
		} else if strings.HasPrefix(m, `$Down`) {
			if strings.HasPrefix(m, `$Down `+remoteAddr) || strings.HasPrefix(m, `$Down all`) {
				if e := conn.Close(); e != nil {
					log.Printf("| %s", e.Error())
					continue
				}
				del(&all, remoteAddr)
				break
			}
		} else {
			_, e := conn.Write([]byte(m))
			if e != nil {
				log.Printf("| %s", e.Error())
				continue
			}
			log.Printf("| %s", m)
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

func delSame(ar *[]string) {
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
	*ar = arr
}

func del(ar *[]string, s string) {
	arr := *ar
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
	if len(*arr) != 0 {
		for i, v := range *arr {
			if i != len(*arr)-1 {
				if v != s {
					continue
				}
				break
			}
			if v != s {
				*arr = append(*arr, s)
			}
		}
		return
	}
	*arr = append(*arr, s)
}
