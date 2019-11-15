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
	Name, Messages string
}
type Display struct {
	Name     []string
	Messages string
	Status   bool
}
type Trans struct {
	TransCall, TransHost []string
	Messages             string
}

func main() {
	flag.Parse()
	scan, calls, ch := bufio.NewScanner(os.Stdin), Call{Name: `[Unkown]`}, make(chan []string)
	reg, _ := regexp.Compile(`[^\s*$]`)

	fmt.Print(`Your name [Default "Unkown"]: `)
	scan.Scan()

	if !reg.MatchString(scan.Text()) {
		calls.Name = fmt.Sprintf("[%s]", scan.Text())
	}

	var trans Trans
	fmt.Print("Transits server host : ")
	scan.Scan()
	if !reg.MatchString(scan.Text()) {
		trans.TransHost = strings.Fields(scan.Text())
		delNearby(&trans.TransHost)
		fmt.Print(`Who do you want to call by transfer: `)
		scan.Scan()
		if !reg.MatchString(scan.Text()) {
			trans.TransCall = strings.Fields(scan.Text())
			delSame(&trans.TransCall)
		}
	}

	var to []string
	fmt.Print(`Who do you want to call by face: `)
	scan.Scan()
	if !reg.MatchString(scan.Text()) {
		to = strings.Fields(scan.Text())
		delSame(&to)
		log.Print(`| Test calling...`)

		for _, v := range to {
			go calls.Test(v)
		}
	}

	for {
		scan.Scan()
		t := scan.Text()
		switch !reg.MatchString(t) {
		case strings.HasPrefix(t, `$Call`):
			ch <- strings.Fields(t)
		case strings.HasPrefix(t, `$TransCall`):
			ch <- strings.Fields(t)
		case strings.HasPrefix(t, `$Name`):
			calls.Name = strings.Fields(t)[1]
		case strings.HasPrefix(t, `$TransHost`):
			ch <- strings.Fields(t)
		default:
			calls.Messages = t
			if len(trans.TransCall) != 0 {
				trans.Messages = t

			}
			for _, v := range to {
				d, e := json.Marshal(calls)
				if e != nil {
					log.Printf("| %s", e.Error())
					continue
				}
				_, e = http.Get(fmt.Sprintf("http://%s/conversation?calls=%s", v, string(d)))
				if e != nil {
					log.Printf("| %s", e.Error())
					continue
				}
			}
		}
	}
}

func serve(host string) {
	var rs Display
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
