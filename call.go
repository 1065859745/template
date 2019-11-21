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

type Call struct {
	port       *string
	udpAddress net.UDPAddr
	messages   string
	all, to    []net.UDPAddr
}

var call = Call{port: flag.String("p", "", "Local serve port, default random.")}

func main() {
	flag.Parse()
	scan := bufio.NewScanner(os.Stdin)
	reg, _ := regexp.Compile(`^\s*$`)

	if *call.port == `` {
		for {
			rand.Seed(time.Now().Unix())
			p := strconv.Itoa(rand.Intn(59000) + 1000)
			*call.port = p
			udpAddr, _ := net.ResolveUDPAddr(`udp`, ":"+*call.port)
			l, e := net.ListenUDP("udp", udpAddr)
			if e != nil {
				continue
			} else {
				l.Close()
			}
			fmt.Printf("Serve at 127.0.0.1:%s\n", *call.port)
			call.udpAddress = *udpAddr
			break
		}
	}
	go serve()

	for {
		scan.Scan()
		t := scan.Text()
		if !reg.MatchString(t) {
			if strings.HasPrefix(t, `$Call`) {
				to := strings.Fields(t)[1:]
				delSame(&to, `127.0.0.1:`+*call.port)
				c := []net.UDPAddr{}
				for _, v := range to {
					addr, e := net.ResolveUDPAddr(`udp`, v)
					if e != nil {
						log.Printf("| %s ", e.Error())
						continue
					}
					c = append(c, *addr)
				}
				call.to = c
				continue
			}
			if strings.HasPrefix(t, `$Clear all`) {
				call.all = []net.UDPAddr{}
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
				call.to = []net.UDPAddr{}
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
	conn, _ := net.ListenUDP("udp", &call.udpAddress)
	for {
		var b [512]byte
		n, addr, e := conn.ReadFromUDP(b[:])
		if e != nil {
			log.Printf("| %s", e.Error())
			return
		}

		m := string(b[:n])
		remoteAddr := fmt.Sprintf("%s:%s", addr.IP, strings.Split(m, ` `)[0])

		if remoteAddr != fmt.Sprintf("127.0.0.1:%s", *call.port) {
			log.Printf("|From| %s:%s", addr.IP, m)
			if e = addUDP(&call.all, remoteAddr); e != nil {
				log.Printf("| %s", e.Error())
			}
		}
	}
}

func (c Call) Send(a []net.UDPAddr) {
	for _, v := range a {
		udpAddr := v.String()
		conn, e := net.Dial("udp", udpAddr)
		if e != nil {
			log.Printf("| %s", e.Error())
			continue
		}

		log.Printf("|To| %s", c.messages)
		conn.Write([]byte(fmt.Sprintf("%s : %s", *c.port, c.messages)))
		conn.Close()

		if udpAddr != `127.0.0.1`+*c.port {
			if e = addUDP(&call.all, udpAddr); e != nil {
				log.Printf("| %s", e.Error())
				return
			}
		}
	}
}

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

func addUDP(u *[]net.UDPAddr, s string) error {
	udpAddr, e := net.ResolveUDPAddr(`udp`, s)
	if e != nil {
		return e
	}
	if len(*u) != 0 {
		for i, v := range *u {
			if i != len(*u)-1 {
				if v.String() != s {
					continue
				}
				break
			}
			if v.String() != s {
				*u = append(*u, *udpAddr)
			}
		}
	} else {
		*u = append(*u, *udpAddr)
	}
	return e
}
