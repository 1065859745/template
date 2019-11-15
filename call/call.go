package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type Call struct {
	Name, Messages, To string
	Transits, ToArr    []string
	useTrans           bool
}
type Result struct {
	Name     []string
	Messages string
	Status   bool
}

func main() {
	flag.Parse()
	scan := bufio.NewScanner(os.Stdin)
	reg, _ := regexp.Compile(`[^\s*$]`)

	fmt.Print(`Your name [Default "Unkown"]: `)
	scan.Scan()

	calls := Call{Name: `[` + scan.Text() + `]`}
	if reg.MatchString(calls.Name) {
		calls.Name = `[Unkown]`
	}

	var transHost string
	fmt.Print("Transits server host : ")
	scan.Scan()
	calls.Transits = delSame(strings.Split(scan.Text(), ` `))
	if !reg.MatchString(calls.Transits[0]) {
		fristT := strings.Split(calls.Transits[0], `:`)
		transHost = calls.Transits[0]
		calls.Transits = calls.Transits[1:]
		if fristT[0] == `127.0.0.1` || fristT[0] == `localhost` {
			if len(fristT) == 2 {
				go serve(calls.Transits[0])
			} else {
				go serve(``)
			}
		}
	} else {
		calls.Transits = []string{}
	}

	fmt.Print(`Who do you want to call: `)
	scan.Scan()
	if !reg.MatchString(scan.Text()) {
		calls.ToArr = strings.Split(scan.Text(), ` `)
		calls.To = strings.Join(calls.ToArr, `,`)
		log.Print(`| Test calling...`)

		if calls.useTrans {
			go calls.Test(transHost)
		} else {
			for _, v := range calls.ToArr {
				go calls.Test(v)
			}
		}
	}

	for {
		scan.Scan()
		text := scan.Text()
		if text == `` {
			continue
		} else if text == `$Call` {

		} else if text == `$Trans` {

		} else if text == `$Name` {
			continue
		} else {
			calls.Messages = scan.Text()
		}
		calls.Println()
		if calls.useTrans {
			calls.Sender(transHost)
		} else {
			for _, v := range calls.ToArr {
				calls.Sender(v)
			}
		}
	}
}

func serve(host string) {
	var rs Result
	rs.Name = []string{`localhost`, ``}
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// query := r.URL.Query()

	})

	err := http.ListenAndServe(host, nil)
	if err != nil {
		rs.Name[1] = `[Error]`
		rs.Messages = err.Error()
		rs.Println()
	}
	return
}

func (cs Call) Sender(host string) {
	d, e := json.Marshal(cs)
	if e != nil {
		log.Printf("| %s", e.Error())
		return
	}

	_, e = http.Get(fmt.Sprintf("http://%s/conversation?calls=%s", host, string(d)))
	if e != nil {
		log.Printf("| %s", e.Error())
	}
	cs.Println()
}

func (cs Call) Test(host string) {
	var rsTest Result
	var testU string
	rsTest.Name = []string{host, ``}
	if cs.useTrans {
		d, e := json.Marshal(cs)
		if e != nil {
			rsTest.Error(e.Error())
			return
		}
		testU = fmt.Sprintf("http://%s/test?calls=%s", host, string(d))
	} else {
		testU = fmt.Sprintf("http://%s/test", host)
	}

	_, e := url.Parse(testU)
	if e != nil {
		rsTest.Error(e.Error())
	}

	r, e := http.Get(testU)
	if e != nil {
		rsTest.Error(e.Error())
	} else if r.StatusCode != 200 {
		rsTest.Status = false
		rsTest.Name[1] = `[Error]`
		rsTest.Messages = fmt.Sprintf("StatusCode %d", r.StatusCode)
	} else if text, e := ioutil.ReadAll(r.Body); e != nil {
		rsTest.Error(e.Error())
	} else {
		rsTest.Status = true
		rsTest.Messages = string(text)
	}
	if rsTest.Status {
		rsTest.Println()
		/* receiver */
		for {
			resp, e := http.Get(testU)
			var respH Call
			if e != nil {
				rsTest.Error(e.Error())
			} else if content, err := ioutil.ReadAll(resp.Body); err != nil {
				if e = json.Unmarshal(content, &respH); err != nil {
					rsTest.Name[1] = respH.Name
					rsTest.Messages = respH.Messages
					rsTest.Println()
				}
			}
			time.Sleep(1 * time.Second)
		}
	} else {
		rsTest.Println()
	}
}

func (rs Result) Error(e string) {
	rs.Name[1] = `[Error]`
	rs.Status = false
	rs.Messages = e
}
func (rs Result) Print() {
	fmt.Printf("%s | %s : %s", time.Now().Format("2006/01/02 15:04:05"), strings.Join(rs.Name, ` `), rs.Messages)
}
func (rs Result) Println() {
	log.Printf("| %s : %s", strings.Join(rs.Name, ` `), rs.Messages)
}
func (cs Call) Println() {
	log.Printf("| %s : %s", cs.Name, cs.Messages)
}
func delSame(arr []string) []string {
a:
	for {
		for i, v := range arr {
			if i != len(arr)-1 {
				if v == arr[i+1] {
					arr = append(arr[:i], arr[i+1:]...)
					continue a
				}
			}
		}
		break
	}
	return arr
}
