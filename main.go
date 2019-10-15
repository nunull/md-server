package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	html := string(markdown.ToHTML(md, nil, nil))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "%s", html)
	})

	port := "8080"

	err = exec.Command("open", fmt.Sprintf("http://localhost:%s", port)).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %s", err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
