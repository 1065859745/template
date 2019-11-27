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

var port, all, messageCh, downCh, answerCh = flag.String("p", "", "Local serve port, default random."), []net.UDPAddr{}, make(chan string, 10), make(chan string, 5), make(chan string)

func main() {
	flag.Parse()
	scan, messages, name := bufio.NewScanner(os.Stdin), ``, `[Unkown]`
	reg, _ := regexp.Compile(`^\s*$`)

	go serve()

	for {
		scan.Scan()
		t := scan.Text()
		if !reg.MatchString(t) {
			tArr := strings.Fields(t)[1:]
			if strings.HasPrefix(t, `$Call`) {
				for _, v := range tArr {
					if udpAddr, e := net.ResolveUDPAddr(`udp`, v); e != nil {
						log.Printf("| %s", e.Error())
					} else {
						if !includeUdp(&all, *udpAddr) {
							all = append(all, *udpAddr)
							go client(udpAddr)
						}
					}
				}
				continue
			}
			if strings.HasPrefix(t, `$Name`) {
				name = strings.Join(strings.Split(t, ` `)[1:], ` `)
				continue
			}
			if strings.HasPrefix(t, `$List`) {
				fmt.Println(`-------------------`)
				for i, v := range all {
					fmt.Println(i+1, v.String())
				}
				fmt.Println(`-------------------`)
				continue
			}
			if strings.HasPrefix(t, `$Answer`) {
				for _, v := range tArr {
					messageCh <- `$Answer ` + v
					<-messageCh
				}
			answer:
				for {
					scan.Scan()
					t = scan.Text()
					if !reg.MatchString(t) {
						t = fmt.Sprintf("%s : %s", name, t)
						answerCh <- t
						<-answerCh
						break answer
					}
				}
				continue
			}
			if strings.HasPrefix(t, `$Down`) {
				for _, v := range tArr {
					delUDPAddr(&all, v)
					downCh <- v
					<-downCh
				}
				continue
			}

			messages = fmt.Sprintf("%s : %s", name, t)
			messageCh <- messages
			<-messageCh
			log.Printf("| %s", messages)
		}
	}
}

func serve() {
	var conn *net.UDPConn
	var udpAddr *net.UDPAddr
	var e error

	/* if *port equal nil,it becomes random*/
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

	/* read */
	var b [512]byte
	for {
		n, addr, _ := conn.ReadFromUDP(b[:])
		log.Printf("|From| %s %s", addr.String(), string(b[:n]))

		if !includeUdp(&all, *addr) {
			all = append(all, *addr)
			go write(conn, addr, true)
		}
	}
}

func client(udpAddr *net.UDPAddr) {
	conn, _ := net.DialUDP("udp", nil, udpAddr)
	addr := udpAddr.String()
	/* write */
	go write(conn, udpAddr, false)

	/* read */
	var b [512]byte
	for {
		n, _, e := conn.ReadFromUDP(b[:])
		if e != nil {
			break
		}
		log.Printf("|From| %s %s", addr, string(b[:n]))
	}
}

func write(conn *net.UDPConn, udpAddr *net.UDPAddr, b bool) {
	var m string
	addr := udpAddr.String()
	defer delUDPAddr(&all, addr)

a:
	for {
		select {
		case m = <-messageCh:
			messageCh <- m
		case m = <-downCh:
			downCh <- m
			if m == addr || m == `all` {
				conn.Close()
				break a
			}
			continue
		}
		if strings.HasPrefix(m, `$Answer`) {
			if strings.HasPrefix(m, "$Answer "+addr) {
				m = <-answerCh
				conn.Write([]byte(m))
				log.Printf("|Answer %s| %s", addr, m)
				answerCh <- m
			}
			continue
		}

		/* udp 的写入数据有两种不同的方式，作为客户机用write(),作为服务端用 writeToUDP() */
		if b {
			conn.WriteToUDP([]byte(m), udpAddr)
		} else {
			conn.Write([]byte(m))
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

func includeUdp(addrArr *[]net.UDPAddr, addr net.UDPAddr) bool {
	b := false
	s := addr.String()
	for _, v := range *addrArr {
		if v.String() == s {
			b = true
			break
		}
	}
	return b
}
