package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	for {
		resp, err := http.Get("http://192.168.0.27:12345/states")
		if err != nil {
			fmt.Printf("Error Get states %s", err)
			os.Exit(1)
		} else {
			defer resp.Body.Close()
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error Processing response's Body %s", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", string(contents))

			fmt.Println("sleep...")
			time.Sleep(time.Duration(10) * time.Second)
		}
	}
}
