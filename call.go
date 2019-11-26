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

var port, all = flag.String("p", "", "Local serve port, default random."), []net.UDPAddr{}
var messageCh, downCh, answerCh chan string

func main() {
	flag.Parse()
	scan, messages, name := bufio.NewScanner(os.Stdin), ``, `[Unkown]`
	reg, _ := regexp.Compile(`^\s*$`)

	go serve()

	for {
		scan.Scan()
		t := scan.Text()
		if !reg.MatchString(t) {
			if strings.HasPrefix(t, `$Call`) {

				continue
			}
			if strings.HasPrefix(t, `$Name`) {

				continue
			}
			if strings.HasPrefix(t, `$List`) {

				continue
			}
			if strings.HasPrefix(t, `$Answer`) {

				continue
			}

			messages = fmt.Sprintf("%s : %s", name, t)
			messageCh <- messages
			<-messageCh
		}
	}
}

func serve() {
	var conn *net.UDPConn
	var udpAddr *net.UDPAddr
	var e error

	if *port == `` {
		for {
			rand.Seed(time.Now().Unix())
			*port = strconv.Itoa(rand.Intn(59000) + 1000)
			udpAddr, _ = net.ResolveUDPAddr(`udp`, ":"+*port)
			conn, e = net.ListenUDP("udp", udpAddr)
			if e != nil {
				continue
			}
			fmt.Printf("Serve at 127.0.0.1:%s\n", *port)
			break
		}
	} else {
		udpAddr, e = net.ResolveUDPAddr(`udp`, ":"+*port)
		if e != nil {
			log.Fatal(e)
		}
		conn, _ = net.ListenUDP("udp", udpAddr)
	}

	addUDP(&all, udpAddr.String())

	go read(conn)

	write(conn)
}

func client(udpAddr *net.UDPAddr) {
	conn, _ := net.DialUDP("udp", nil, udpAddr)

	addUDP(&all, udpAddr.String())

	go read(conn)

	write(conn)
}

func read(conn *net.UDPConn) {
	remoteAddr := conn.RemoteAddr().String()
	for {
		var b [512]byte
		n, addr, e := conn.ReadFromUDP(b[:])
		if e != nil {
			downCh <- `$Down ` + remoteAddr
			conn.Close()
			delUDPAddr(&all, remoteAddr)
			<-downCh
			return
		}

		m := string(b[:n])
		log.Printf("|From| %s %s", addr.String(), m)
		log.Printf("| %s", e.Error())
	}
}

func write(conn *net.UDPConn) {
	var m string
	remoteAddr := conn.RemoteAddr().String()

	for {
		select {
		case m = <-messageCh:
			messageCh <- m
		case m = <-downCh:
			downCh <- m
		}

		if strings.HasPrefix(m, "$Answer "+remoteAddr) {
			m = <-answerCh
			answerCh <- m
		} else if strings.HasPrefix(m, "$Down "+remoteAddr) || strings.HasPrefix(m, "$Down all") {
			conn.Close()
			break
		}

		_, e := conn.Write([]byte(m))
		if e != nil {
			log.Printf("| %s", e.Error())
			break
		}
	}
}

func delSame(arr *[]string) {
	ar := *arr
a:
	for {
		for i, n := range ar {
			if i != len(ar)-1 {
				for j, m := range ar[i+1:] {
					if n != m {
						continue
					} else {
						ar = append(ar[:j+i+1], ar[j+i+2:]...)
						continue a
					}
				}
			}
		}
		break
	}
	*arr = ar
}

func delUDPAddr(arr *[]net.UDPAddr, s string) {
	ar := *arr
	for i, v := range ar {
		if v.String() != s {
			continue
		}
		ar = append(ar[:i], ar[i+1:]...)
		break
	}
	*arr = ar
}

func addUDP(u *[]net.UDPAddr, s string) {
	udpAddr, _ := net.ResolveUDPAddr(`udp`, s)

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
}
