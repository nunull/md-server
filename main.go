package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gomarkdown/markdown"
)

var rootDirname string
var address string

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <dirname>\n", os.Args[0])
		os.Exit(1)
	}
	rootDirname = os.Args[1]

	http.HandleFunc("/", handle)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address = "http://localhost:" + port
	if err := exec.Command("open", address).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %s", err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	serveFileOrDir(w, r.URL.Path)
}

func serveFileOrDir(w http.ResponseWriter, filename string) {
	fi, err := os.Stat(rootDirname + filename)
	if err != nil {
		// TODO respond
		fmt.Println(err)
		return
	}

	var html string
	if fi.Mode().IsRegular() {
		html, err = getFileHtml(w, filename)
		if err != nil {
			// TODO respond
			fmt.Println(err)
			return
		}
	} else if fi.Mode().IsDir() {
		html, err = getDirHtml(w, filename)
		if err != nil {
			// TODO respond
			fmt.Println(err)
			return
		}
	}

	fullHtml := "<style>* { font-family: sans-serif; } code { font-family: monospace; }</style>"
	fullHtml += "<header>"
	fullHtml += "<a href=\"/\">index</a>"
	pathPrefix := ""
	for _, item := range strings.Split(filename, "/")[1:] {
		pathPrefix += "/" + item
		fullHtml += "/<a href=\"" + pathPrefix + "\">" + item + "</a>"
	}
	fullHtml += "</header>"
	fullHtml += html

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%s", fullHtml)
}

func getFileHtml(w http.ResponseWriter, filename string) (string, error) {
	md, err := ioutil.ReadFile(rootDirname + filename)
	if err != nil {
		return "", err
	}

	return string(markdown.ToHTML(md, nil, nil)), nil
}

func getDirHtml(w http.ResponseWriter, dirname string) (string, error) {
	files, err := ioutil.ReadDir(rootDirname + dirname)
	if err != nil {
		return "", err
	}

	html := "<ul>"
	for _, file := range files {
		html += "<li><a href=\"" + address + "/" + dirname + "/" + file.Name() + "\">" + file.Name() + "</a>"
	}
	html += "</ul>"

	return html, nil
}
