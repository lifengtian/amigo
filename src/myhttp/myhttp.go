
package main

import (
	"io"
	"net/http"
	"log"
	"time"
	"fmt"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func ShowStates(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "It's %s\n", time.Now() )
}

func SetPassword (w http.ResponseWriter, req *http.Request) {
	req.SetBasicAuth("test", "garyowen")	
}
	
func main() {
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/states", ShowStates )
	http.HandleFunc("/setpwd", SetPassword )
	err := http.ListenAndServe(":12345", nil)
	
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}