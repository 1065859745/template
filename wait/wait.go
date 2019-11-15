package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var port = flag.String("p", "80", "Start port for waiting server.")

type Call struct {
	Name, Messages, To string
	Transits, ToArr    []string
}

type Response struct {
	Name     []string
	Messages string
	Status   bool
}

var display Response
var call Call

func main() {
	flag.Parse()
	scan := bufio.NewScanner(os.Stdin)

	ch := make(chan string, 20)
	reg, _ := regexp.Compile(`[^\s*$]`)

	fmt.Print(`Your name [Default "Unkown"]: `)
	scan.Scan()
	if scan.Text() != `` {
		call.Name = fmt.Sprintf("[%s]", scan.Text())
	} else {
		call.Name = `[Unkown]`
	}
	log.Printf("| Serve will start at :%s...", *port)
	go serve(*port, ch)
	for {
		scan.Scan()
		if !reg.MatchString(scan.Text()) {
			call.Messages = scan.Text()
			call.Println()
			ch <- display.Messages
		}
	}
}
func serve(port string, c chan string) {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		back := `Ok`
		display = Response{Name: []string{r.Host, ``}, Messages: `Ok`}
		log.Printf("| %s connecting...\n", display.Name[0])

		query := r.URL.Query()
		if calls, ok := query["calls"]; ok {
			var cs Call
			if e := json.Unmarshal([]byte(calls[0]), &cs); e != nil {
				back = e.Error()
			} else if len(cs.Transits) != 0 {
				/* Transfer */
				if r, e := http.Get(cs.Transits[0]); e != nil {
					back = e.Error()
				} else if t, e := ioutil.ReadAll(r.Body); e != nil {
					back = e.Error()
				} else {
					back = string(t)
				}
			}
		}
		fmt.Fprint(w, back)
	})

	http.HandleFunc("/conversation", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		display.Name = []string{r.Host, `[Error]`}
		back := `Ok`
		if calls, ok := query["calls"]; ok {
			/* parser come message */
			var cs Call
			if e := json.Unmarshal([]byte(calls[0]), &cs); e != nil {
				back = e.Error()
			} else if len(cs.Transits) != 0 {
				/* Transfer */
				if r, e := http.Get(cs.Transits[0]); e != nil {
					back = e.Error()
				} else if t, e := ioutil.ReadAll(r.Body); e != nil {
					back = e.Error()
				} else {
					back = string(t)
				}
			} else {
				/* face to face */
				display.Name[1] = cs.Name
				display.Messages = cs.Messages
				display.Println()
				call.Messages = <-c
				d, _ := json.Marshal(call)
				back = string(d)
			}
		}
		fmt.Fprint(w, back)
	})
	resp := Response{Name: []string{`localhost`, ``}}
	err := http.ListenAndServe(`:`+port, nil)
	if err != nil {
		resp.Name[1] = `[Error]`
		resp.Messages = err.Error()
		resp.Println()
		return
	}
}
func (rs Response) Println() {
	log.Printf("| %s : %s", strings.Join(rs.Name, ` `), rs.Messages)
}
func (cs Call) Println() {
	log.Printf("| %s : %s", cs.Name, cs.Messages)
}
