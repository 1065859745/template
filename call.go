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
var chAnswer, chDown chan string

func main() {
	flag.Parse()
	var scan = bufio.NewScanner(os.Stdin)
	reg, _ := regexp.Compile(`^\s*$`)

	switch *t {
	case -1:
		// serve()
		fmt.Println(`-1`)
	case 1:
		fmt.Println(`1`)
		// client()
	default:
		go serve()

		var t string
		var tArr []string
		for {
			scan.Scan()
			t = scan.Text()
			if !reg.MatchString(t) {
				tArr = strings.Fields(t)[1:]
				if strings.HasPrefix(t, `$Name`) {
					nameArr := strings.Split(t, ` `)[1:]
					name = `[` + strings.Join(nameArr, ` `) + `]`
					continue
				}
				if strings.HasPrefix(t, `$Call`) {
					delSame(&tArr)
					for _, v := range tArr {
						go client(v)
					}
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
				if strings.HasPrefix(t, `$Answer`) {
					for _, v := range tArr {
						chStr <- `$Answer ` + v
						<-chStr
					}
				answer:
					for {
						scan.Scan()
						t = scan.Text()
						if !reg.MatchString(t) {
							t = fmt.Sprintf("%s : %s", name, t)
							chAnswer <- t
							<-chAnswer
							break answer
						}
					}
					continue
				}
				if strings.HasPrefix(t, `$Down`) {
					for _, v := range tArr {
						chStr <- `$Down ` + v
						<-chStr
					}
					continue
				}
				if strings.HasPrefix(t, `$Port`) {
					fmt.Println(`------------`)
					fmt.Printf(":%s\n", *port)
					fmt.Println(`------------`)
					continue
				}
				t = fmt.Sprintf("%s : %s", name, t)
				chStr <- t
				<-chStr
				log.Printf("| %s", t)
			}
		}
	}
}

func serve() {
	var l net.Listener
	var e error
	if *port == `` {
		for {
			rand.Seed(time.Now().Unix())
			*port = strconv.Itoa(rand.Intn(59000) + 1000)
			l, e = net.Listen("tcp", ":"+*port)
			if e != nil {
				continue
			}
			fmt.Printf("Serve at 127.0.0.1:%s\n", *port)
			break
		}
	} else if l, e = net.Listen("tcp4", ":"+*port); e != nil {
		log.Fatal(e)
	}

	for {
		conn, e := l.Accept()
		if e != nil {
			log.Printf("| %s", e.Error())
			continue
		}
		remoteAddr := conn.RemoteAddr().String()
		update(&all, remoteAddr)
		log.Printf(`Connected to %s successful.`, remoteAddr)
		// read
		go tcpRead(conn)

		// write
		tcpWrite(conn)
	}
}

func client(addr string) {
	conn, e := net.Dial("tcp4", addr)
	if e != nil {
		log.Printf("| %s", e.Error())
		return
	}
	update(&all, addr)
	log.Printf(`Connected to %s successful.`, addr)

	// read
	go tcpRead(conn)

	// write
	tcpWrite(conn)
}

func tcpRead(conn net.Conn) {
	var b [512]byte
	remoteAddr := conn.RemoteAddr().String()
	defer del(&all, remoteAddr)
	for {
		n, e := conn.Read(b[:])
		if e != nil {
			chDown <- `$Down ` + remoteAddr
			log.Printf("|Miss| %s", remoteAddr)
			<-chDown
			break
		}
		log.Printf("|From| %s %s", remoteAddr, string(b[:n]))
	}
}

func tcpWrite(conn net.Conn) {
	var m string
	remoteAddr := conn.RemoteAddr().String()

	for {
		select {
		case m = <-chStr:
			chStr <- m
		case m = <-chDown:
			chDown <- m
		}

		if strings.HasPrefix(m, `$Answer`) {
			if strings.Fields(m)[1] == remoteAddr {
				m = <-chAnswer
				if !strings.HasPrefix(m, `$`) {
					log.Printf("|Answer %s| %s", remoteAddr, m)
				} else {
					continue
				}
				chAnswer <- m
			}
		} else if strings.HasPrefix(m, `$Down`) {
			if strings.HasPrefix(m, `$Down `+remoteAddr) || strings.HasPrefix(m, `$Down all`) {
				conn.Close()
				break
			}
			continue
		}

		_, e := conn.Write([]byte(m))
		if e != nil {
			log.Printf("| %s", e.Error())
			continue
		}
	}
}

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
func includes(arr []string, s string) bool {
	b := false
	for _, v := range arr {
		if v == s {
			b = true
			break
		}
	}
	return b
}
