package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// package main

// import (
// 	"bytes"
// 	"flag"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"mime/multipart"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// )

// var transfer, end = flag.String("h", "", "Transfer ip"), flag.String("e", "", "End ip")

// func main() {
// 	flag.Parse()
// 	if *transfer == "" {
// 		fmt.Print("Must enter transfer ip")
// 		return
// 	}
// 	if *end == "" {
// 		fmt.Print("Must enter end ip")
// 		return
// 	}
// 	bodyBuffer := &bytes.Buffer{}
// 	bodyWriter := multipart.NewWriter(bodyBuffer)
// 	for _, v := range flag.Args() {
// 		fileWriter, _ := bodyWriter.CreateFormFile("files", filepath.Base(v))
// 		file, _ := os.Open(v)
// 		io.Copy(fileWriter, file)
// 		file.Close()
// 	}
// 	contentType := bodyWriter.FormDataContentType()
// 	bodyWriter.Close()

// 	resp, _ := http.Post("http://"+*transfer, contentType, bodyBuffer)

// 	resp_body, _ := ioutil.ReadAll(resp.Body)

// 	log.Println(resp.Status)
// 	log.Println(string(resp_body))
// 	defer resp.Body.Close()
// }

var port = flag.String("p", "80", "Server start port")

func main() {
	flag.Parse()
	param := flag.Args()
	var mym string
	fmt.Println("Enter what you want to say to remote ip")
	fmt.Print("KaoYaDian: ")
	fmt.Scanln(&mym)
	resp, e := http.Get("http://" + param[0] + "/?to=\"" + param[1] + "\";messages=\"" + mym + "\"")
	if e != nil {
		log.Fatal(e.Error())
	}
	resp.
		log.Fatal(http.ListenAndServe(":"+*port, nil))
}
